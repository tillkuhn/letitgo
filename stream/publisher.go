package stream

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go/sasl/plain"

	"github.com/kelseyhightower/envconfig"
	"github.com/segmentio/kafka-go"
)

type KafkaPublisher struct {
	config KafkaConfig
	logger zerolog.Logger
}

// NewKafkaPublisher returns a new configured and ready-to-use KafkaPublisher
func NewKafkaPublisher() (*KafkaPublisher, error) {
	logger := log.With().Str("logger", "kafka-publisher").Logger()
	var config KafkaConfig
	if err := envconfig.Process(envconfigPrefix, &config); err != nil {
		log.Error().Msgf("Cannot parse envconfig: %v", err)
		return nil, err
	}
	if !config.Enabled {
		return nil, errors.New("kafka is disabled by config")
	}
	p := &KafkaPublisher{config: config, logger: logger}
	logger.Info().Msgf("NewKafkaPublisher client config broker=%s topic=%s", p.config.Servers, config.Topic)
	return p, nil
}

// WriteMessage ✉️ me a letter, will you?
func (p *KafkaPublisher) WriteMessage(ctx context.Context, key string, message []byte) error {
	return p.WriteMessageToTopic(ctx, p.config.Topic, key, message)
}

// WriteMessageToTopic ✉️ me a letter, will you? You can even pick a topic :-)
func (p *KafkaPublisher) WriteMessageToTopic(ctx context.Context, topic string, key string, message []byte) error {
	w := &kafka.Writer{
		Addr:  kafka.TCP(p.config.Servers), // can take multiple servers todo: split arg if > 1 brokers
		Topic: topic,
		// balancer to determine which partition to route messages to
		// kafka.Hash is similar to Sarama Algorithm to pick the partition, Murmur2Balancer is similar to java impl
		Balancer:     &kafka.Murmur2Balancer{},
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		RequiredAcks: kafka.RequireOne, //  wait for the leader to acknowledge the writes
		Transport:    p.transport(),
		Logger:       LoggerWrapper{p.logger},
	}

	// Let's deliver our important message
	err := w.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(key),
			Value: message,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to write messages to %s: %v", topic, err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to close writer for %s: %v", topic, err)
	}
	p.logger.Info().Msgf("Successfully pushed %d bytes to topic %s", len(message), topic)
	return nil
}

// WriteStringMessage Convenience method to pass a string as message payload
func (p *KafkaPublisher) WriteStringMessage(ctx context.Context, key string, message string) error {
	return p.WriteMessage(ctx, key, []byte(message))
}

// WriteJSONMessage Convenience method to pass a string as message payload
func (p *KafkaPublisher) WriteJSONMessage(ctx context.Context, key string, message interface{}) error {
	responseBytes, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		return err
	}
	return p.WriteMessage(ctx, key, responseBytes)
}

// CreateTopic creates a new Topic to create topics when auto.create.topics.enable='false'
// if a topic does not exist, you'll be writing messages wil fail with
// [3] Unknown Topic Or Partition: the request is for a topic or partition that does not exist on this broker
func (p *KafkaPublisher) CreateTopic(topic string, partitions int) error {
	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		TLS:           &tls.Config{MinVersion: tls.VersionTLS12},
		SASLMechanism: plain.Mechanism{Username: p.config.APIKey, Password: p.config.APISecret},
	}

	conn, err := dialer.Dial("tcp", p.config.Servers)
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}
	var controllerConn *kafka.Conn
	controllerConn, err = dialer.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	p.logger.Debug().Msgf("About to create topic %s with %d partitions", topic, partitions)
	topicConfigs := []kafka.TopicConfig{
		{
			Topic:         topic,
			NumPartitions: partitions,
			// for some reason we get "the request parameters do not satisfy the configured policy if this is too low
			ReplicationFactor: 3,
		},
	}

	if err = controllerConn.CreateTopics(topicConfigs...); err != nil {
		return err
	}
	p.logger.Debug().Msgf("Topic %s with %d partitions successfully created", topic, partitions)
	return nil
}

// Config returns a copy of the KafkaConfig (so modifications from the client have no effect)
func (p *KafkaPublisher) Config() KafkaConfig {
	return p.config
}

// transport returns an implementation of the Transport Round tripper interface with SASL and TLS enabled
func (p *KafkaPublisher) transport() *kafka.Transport {
	return &kafka.Transport{
		ClientID: p.config.ClientID,
		SASL: &plain.Mechanism{ // plain doesn't mean not-encrypted, in just means we don't use SCRAM
			Username: p.config.APIKey,
			Password: p.config.APISecret,
		},
		TLS: &tls.Config{MinVersion: tls.VersionTLS12}, //  If the TLS field is nil, it will not connect with TLS.
	}
}

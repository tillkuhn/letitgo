package stream

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

const envconfigPrefix = "kafka"

// KafkaConfig Kafka Context params populated by envconfig in NewConsumer...()
type KafkaConfig struct {
	APIKey    string `required:"false" default:"" desc:"Kafka API Key Key (user)"  split_words:"true"`
	APISecret string `required:"false" default:"" desc:"Kafka API Secret (password)" split_words:"true"`
	ClientID  string `required:"false" default:"troxyClient" desc:"Client Id for Message Producer" split_words:"true"`
	Debug     bool   `default:"false" desc:"Debug mode, registers logger for kafka packages" split_words:"true"`
	Enabled   bool   `default:"false" desc:"Toggle Kafka Integration on / off (default: false)" split_words:"true"`
	GroupID   string `required:"false" default:"troxyLocal" desc:"Used  as default id for KafkaConsumerGroups" split_words:"true"`
	Servers   string `required:"false" default:"localhost:9092" desc:"Kafka Bootstrap servers" split_words:"true"`
	Topic     string `required:"false" default:"" desc:"Name of Kafka Topic to push tu" split_words:"true"`
}

type KafkaConsumer struct {
	CancelFunc context.CancelFunc

	consumerFunc func(message kafka.Message)
	config       KafkaConfig
	logger       zerolog.Logger
	doneChan     chan struct{}
}

// NewConsumer returns a new configured and ready-to-use KafkaPublisher
// See https://github.com/segmentio/kafka-go#reader-
func NewConsumer(consumerFunc func(message kafka.Message)) (*KafkaConsumer, error) {
	logger := log.With().Str("logger", "kafka-consumer🐛").Logger()
	var config KafkaConfig
	if err := envconfig.Process(envconfigPrefix, &config); err != nil {
		logger.Error().Msgf("Cannot parse envconfig: %v", err)
		return nil, err
	}
	if config.Topic == "" {
		return nil, errors.New("cannot create consumer with empty topic")
	}
	k := &KafkaConsumer{
		consumerFunc: consumerFunc,
		config:       config,
		logger:       logger,
	}
	logger.Debug().Msgf("New KafkaConsumer client config broker=%s consumerGroupId=%s", k.config.Servers, k.config.GroupID)
	return k, nil
}

// Start starts consumeMessage loop in a separate go routine, using a cancel context
// that allows us to shut down things with Stop
func (k *KafkaConsumer) Start() {
	if !k.config.Enabled {
		k.logger.Warn().Msg("Kafka is disabled, skip start kafkaesk")
		return
	}
	errChan := make(chan error, 10)
	doneChan := make(chan struct{})

	// WithCancel returns a copy of parent with a new Done channel. The returned
	// context's Done channel is closed when the returned cancel function is called
	// or when the parent context's Done channel is closed, whichever happens first.
	ctx, cancelConsumer := context.WithCancel(context.Background())
	k.doneChan = doneChan
	k.CancelFunc = cancelConsumer
	k.logger.Info().Msgf("Starting Kafka Consumer client=%s", k.config.ClientID)
	go func() {
		// incoming messages received by the kafkaesk loop will be passed to ConsumerFunc
		errChan <- k.consumeMessages(ctx, doneChan)
	}()
}

// Stop calls the cancel function, and waits for the done channel
func (k *KafkaConsumer) Stop() {
	if !k.config.Enabled {
		k.logger.Warn().Msg("kafka is disabled, skip stop kafkaesk")
		return
	}
	if k.CancelFunc == nil {
		k.logger.Warn().Msg("No cancel func available, KafkaConsumer has not been started")
		return
	}
	k.logger.Debug().Msg("Stopping KafkaConsumer via cancelFunc")
	k.CancelFunc()
	// wait for kafkaesk done channel to close
	<-k.doneChan
	k.logger.Info().Msgf("All channels are closed, goodbye ✌️")
}

// consumeMessages uses kafka-go Reader which automatically handles reconnections and offset management,
// and exposes an API that supports asynchronous cancellations and timeouts using Go contexts.
// See https://github.com/segmentio/kafka-go#reader-
// and this nice tutorial https://www.sohamkamani.com/golang/working-with-kafka/
// doneChan chan<- struct{}
func (k *KafkaConsumer) consumeMessages(ctx context.Context, doneChan chan<- struct{}) error {
	k.logger.Info().Msgf("Let's consume some yummy Kafka Messages on topic=%s", k.config.Topic)
	dialer := &kafka.Dialer{
		Timeout: 3 * time.Second,
		TLS:     &tls.Config{MinVersion: tls.VersionTLS12},
		SASLMechanism: plain.Mechanism{
			Username: k.config.APIKey,
			Password: k.config.APISecret,
		},
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{k.config.Servers},
		GroupID: k.config.GroupID,
		Topic:   k.config.Topic,
		// kafkaesk polls the cluster to check if there is any new data on the topic for the my-group kafkaesk ID,
		// the cluster will only respond if there are at least 10 new bytes of information to send.
		MinBytes: 10,
		MaxBytes: 10e6, // 10MB
		Dialer:   dialer,
		// LastOffset  int64 = -1 // The most recent offset available for a partition.
		// FirstOffset int64 = -2 // The least recent offset available for a partition.
		StartOffset: kafka.LastOffset,
		Logger:      LoggerWrapper{delegate: k.logger},
	})

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			k.logger.Error().Msgf("error on read: %v", err)
			break
		}
		k.consumerFunc(msg)
		// msgChan <- msg
	}

	k.logger.Debug().Msgf("Exited from ReadMessageLoop, closing reader pid=%d", os.Getpid())
	if err := r.Close(); err != nil {
		return fmt.Errorf("failed to close reader: %w", err)
	}
	k.logger.Info().Msgf("Reader successfully closed, closing done channel pid=%d", os.Getpid())
	close(doneChan) // Stop method will wait (block) for this channel before returning
	return nil
}

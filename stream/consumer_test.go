package stream

import (
	"os"
	"strings"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	_, err := testInstance()
	assert.Error(t, err) // topic is mandatory
	_ = os.Setenv(strings.ToUpper(envconfigPrefix)+"_TOPIC", "hase")
	k2, err2 := testInstance()
	assert.NotNil(t, k2)
	assert.NoError(t, err2)
	_ = os.Unsetenv(strings.ToUpper(envconfigPrefix) + "_TOPIC")
	k2.Start() // safe since it's disabled by default
	k2.Stop()  // safe since it's disabled by default
}

func testInstance() (*KafkaConsumer, error) {
	cf := func(msg kafka.Message) {}
	k, err := NewConsumer(cf)
	return k, err
}

// Package kafka provides producer and consumer to work with kafka topics
package kafka

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIntegrationNewConsumer will be passed only if kafka broker is started on localhost:9092
func TestIntegrationNewConsumer(t *testing.T) {
	config := &Config{
		Host:            "localhost",
		Port:            "9092",
		Version:         "2.4.1",
		ConsumerTopic:   "testTopic",
		ConsumerGroupID: "1",
	}
	_, err := NewConsumer(config)

	assert.Nil(t, err)
}

func TestNewConsumerIncorrectVersion(t *testing.T) {
	config := &Config{
		Host:            "localhost",
		Port:            "9092",
		Version:         "11111",
		ConsumerGroupID: "1",
	}

	consumer, err := NewConsumer(config)
	assert.EqualError(t, err, "invalid version `11111`")
	assert.Nil(t, consumer)
}

func TestNewConsumerIncorrectHost(t *testing.T) {
	config := &Config{
		Host:            "localghost",
		Port:            "0",
		Version:         "2.4.1",
		ConsumerGroupID: "1",
	}

	consumer, err := NewConsumer(config)
	assert.EqualError(t, err, "kafka: client has run out of available brokers to talk to (Is your cluster reachable?)")
	assert.Nil(t, consumer)
}

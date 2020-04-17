// Package kafka provides producer and consumer to work with kafka topics
package kafka

import (
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	logger = logrus.New()
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
	consumer, err := NewConsumer(config, logger)
	assert.NotNil(t, consumer)
	assert.NotNil(t, consumer.closeChan)
	assert.Equal(t, config, consumer.Config)
	assert.NoError(t, err)
}

func TestNewConsumerIncorrectVersion(t *testing.T) {
	config := &Config{
		Host:            "localhost",
		Port:            "9092",
		Version:         "11111",
		ConsumerGroupID: "1",
	}

	consumer, err := NewConsumer(config, logger)
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

	consumer, err := NewConsumer(config, logger)
	assert.EqualError(t, err, "kafka: client has run out of available brokers to talk to (Is your cluster reachable?)")
	assert.Nil(t, consumer)
}

func TestConsumerClose(t *testing.T) {
	consumer := &Consumer{
		closeChan: make(chan bool),
		once:      sync.Once{},
	}
	err := consumer.Close()
	assert.NoError(t, err)

	select {
	case _, ok := <-consumer.closeChan:
		assert.False(t, ok)
	default:
		t.Error("Channel is not closed")
	}
}

func TestConsumerDoubleClose(t *testing.T) {
	consumer := &Consumer{
		closeChan: make(chan bool),
		once:      sync.Once{},
	}
	err := consumer.Close()
	assert.NoError(t, err)
	err = consumer.Close()
	assert.NoError(t, err)

	select {
	case _, ok := <-consumer.closeChan:
		assert.False(t, ok)
	default:
		t.Error("Channel is not closed")
	}
}

func TestGetKafkaAddr(t *testing.T) {
	config := &Config{
		Host: "localhost",
		Port: "80",
	}
	expectedAddr := []string{"localhost:80"}
	addr := getKafkaAddr(config)
	assert.Equal(t, expectedAddr, addr)
}

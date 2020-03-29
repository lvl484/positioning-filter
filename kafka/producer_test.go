// Package kafka provides producer and consumer to work with kafka topics
package kafka

import (
	"errors"
	"math"
	"testing"

	saramaMocks "github.com/Shopify/sarama/mocks"
	"github.com/google/uuid"
	"github.com/lvl484/positioning-filter/position"
	"github.com/stretchr/testify/assert"
)

func TestNewProducerFail(t *testing.T) {
	config := &Config{
		Host: "somehost",
		Port: "1",
	}

	got, err := NewProducer(config)

	assert.NotNil(t, err)
	assert.Nil(t, got)
}

// TestNewProducerIntegrationSuccess will be passed only if kafka broker is started on localhost:9092
func TestNewProducerIntegrationSuccess(t *testing.T) {
	config := &Config{
		Host: "localhost",
		Port: "9092",
	}
	got, err := NewProducer(config)

	assert.Nil(t, err)
	assert.NotNil(t, got)
}

func TestProducerClose(t *testing.T) {
	kafkaProducer := saramaMocks.NewSyncProducer(t, nil)
	producer := producer{
		KafkaProducer: kafkaProducer,
	}

	err := producer.Close()
	assert.Nil(t, err)
}

func TestProducerProduceSuccess(t *testing.T) {
	kafkaProducer := saramaMocks.NewSyncProducer(t, nil)
	config := &Config{
		ProducerTopic: "Topic",
	}
	producer := producer{
		KafkaProducer: kafkaProducer,
		Config:        config,
	}

	kafkaProducer.ExpectSendMessageAndSucceed()

	var p position.Position
	err := producer.Produce(p)
	assert.Nil(t, err)
}

func TestProducerProduceFail(t *testing.T) {
	kafkaProducer := saramaMocks.NewSyncProducer(t, nil)
	config := &Config{
		ProducerTopic: "Topic",
	}
	producer := producer{
		KafkaProducer: kafkaProducer,
		Config:        config,
	}

	testError := errors.New("Test error for Producer")

	kafkaProducer.ExpectSendMessageAndFail(testError)

	var p position.Position
	err := producer.Produce(p)
	assert.EqualError(t, err, "Test error for Producer")
}

func TestProducerProduceFailEncode(t *testing.T) {
	kafkaProducer := saramaMocks.NewSyncProducer(t, nil)
	config := &Config{
		ProducerTopic: "Topic",
	}
	producer := producer{
		KafkaProducer: kafkaProducer,
		Config:        config,
	}

	lat := float32(math.Inf(1))
	p := position.Position{
		UserID:   uuid.New(),
		Latitude: lat,
	}

	err := producer.Produce(p)
	assert.NotNil(t, err)
}

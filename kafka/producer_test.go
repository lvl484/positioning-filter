// Package kafka provides producer and consumer to work with kafka topics
package kafka

import (
	"encoding/json"
	"errors"
	"math"
	"testing"
	"time"

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

	var latitude float32 = 1.666

	var longitude float32 = 1.333

	p := position.Position{
		UserID:    uuid.New(),
		Latitude:  latitude,
		Longitude: longitude,
		Timestamp: time.Now(),
		Arrival:   time.Now().Add(time.Second),
	}

	b1, err := json.Marshal(p)
	assert.Nil(t, err)

	checker := func(b2 []byte) error {
		assert.Equal(t, b1, b2)
		return nil
	}

	kafkaProducer.ExpectSendMessageWithCheckerFunctionAndSucceed(checker)

	err = producer.Produce(p)
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

	var latitude float32 = 1.666

	var longitude float32 = 1.333

	p := position.Position{
		UserID:    uuid.New(),
		Latitude:  latitude,
		Longitude: longitude,
		Timestamp: time.Now(),
		Arrival:   time.Now().Add(time.Second),
	}

	b1, err := json.Marshal(p)
	assert.Nil(t, err)

	checker := func(b2 []byte) error {
		assert.Equal(t, b1, b2)
		return nil
	}

	kafkaProducer.ExpectSendMessageWithCheckerFunctionAndFail(checker, testError)

	err = producer.Produce(p)
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

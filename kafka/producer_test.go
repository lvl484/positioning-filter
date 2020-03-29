// Package kafka provides producer and consumer to work with kafka topics
package kafka

import (
	"errors"
	"math"
	"testing"

	saramaMocks "github.com/Shopify/sarama/mocks"
	"github.com/google/uuid"
	"github.com/lvl484/positioning-filter/position"
)

func TestNewProducerFail(t *testing.T) {
	config := &Config{
		Host: "somehost",
		Port: "1",
	}

	got, err := NewProducer(config)
	if err == nil {
		t.Errorf("NewProducer() got nil, want error")
		return
	}

	if got != nil {
		t.Errorf("NewProducer() got Producer, want nil")
	}
}

// TestNewProducerIntegrationSuccess will be passed only if kafka broker is started on localhost:9092
func TestNewProducerIntegrationSuccess(t *testing.T) {
	config := &Config{
		Host: "localhost",
		Port: "9092",
	}
	got, err := NewProducer(config)

	if err != nil {
		t.Errorf("NewProducer() got %v, want nil", err)
	}

	if got == nil {
		t.Errorf("NewProducer() got nil, want Producer")
	}
}

func TestProducerClose(t *testing.T) {
	kafkaProducer := saramaMocks.NewSyncProducer(t, nil)
	producer := producer{
		KafkaProducer: kafkaProducer,
	}

	if err := producer.Close(); err != nil {
		t.Errorf("Close() got %v, want nil", err)
	}
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
	if err := producer.Produce(p); err != nil {
		t.Errorf("Producer.Produce() got %v, want nil", err)
	}
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
	if err := producer.Produce(p); err == nil {
		t.Errorf("Producer.Produce() got nil, want %v", testError)
	}
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

	testError := errors.New("Test error for json.Marshal in Produce()")

	lat := float32(math.Inf(1))
	p := position.Position{
		UserID:   uuid.New(),
		Latitude: lat,
	}

	if err := producer.Produce(p); err == nil {
		t.Errorf("Producer.Produce() got nil, want %v", testError)
	}
}

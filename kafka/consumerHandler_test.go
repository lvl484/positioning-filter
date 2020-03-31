// Package kafka provides producer and consumer to work with kafka topics
package kafka

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	saramaMocks "github.com/Shopify/sarama/mocks"
	"github.com/google/uuid"
	"github.com/lvl484/positioning-filter/matcher"
	"github.com/lvl484/positioning-filter/position"
	"github.com/lvl484/positioning-filter/repository"
	"github.com/stretchr/testify/assert"
)

func TestNewConsumerGroupHandler(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.Nil(t, err)

	filters := repository.NewFiltersRepository(db)
	matcher := matcher.NewMatcher(filters)
	saramaProducer := saramaMocks.NewSyncProducer(t, nil)
	producer := &producer{
		KafkaProducer: saramaProducer,
	}
	consumerGroupHandler := newConsumerGroupHandler(matcher, producer)
	assert.NotNil(t, consumerGroupHandler.controller.matcher)
	assert.NotNil(t, consumerGroupHandler.controller.producer)
}

func TestMessageControllerHandleMessageSuccess(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.Nil(t, err)

	filters := repository.NewFiltersRepository(db)
	matcher := matcher.NewMatcher(filters)
	saramaProducer := saramaMocks.NewSyncProducer(t, nil)
	producer := &producer{
		KafkaProducer: saramaProducer,
		Config: &Config{
			ProducerTopic: "filtered-positions",
		},
	}
	consumerGroupHandler := newConsumerGroupHandler(matcher, producer)

	var latitude float32 = 1.666

	var longitude float32 = 1.333

	pos := position.Position{
		UserID:    uuid.New(),
		Latitude:  latitude,
		Longitude: longitude,
		Timestamp: time.Now(),
		Arrival:   time.Now().Add(time.Second),
	}

	b1, err := json.Marshal(pos)
	assert.Nil(t, err)

	checker := func(b2 []byte) error {
		assert.Equal(t, b1, b2)
		return nil
	}

	saramaProducer.ExpectSendMessageWithCheckerFunctionAndSucceed(checker)

	err = consumerGroupHandler.controller.handleMessage(pos)
	assert.Nil(t, err)
}

func TestMessageControllerHandleMessageFail(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.Nil(t, err)

	filters := repository.NewFiltersRepository(db)
	matcher := matcher.NewMatcher(filters)
	saramaProducer := saramaMocks.NewSyncProducer(t, nil)
	producer := &producer{
		KafkaProducer: saramaProducer,
		Config: &Config{
			ProducerTopic: "filtered-positions",
		},
	}
	consumerGroupHandler := newConsumerGroupHandler(matcher, producer)

	var latitude float32 = 1.666

	var longitude float32 = 1.333

	pos := position.Position{
		UserID:    uuid.New(),
		Latitude:  latitude,
		Longitude: longitude,
		Timestamp: time.Now(),
		Arrival:   time.Now().Add(time.Second),
	}

	b1, err := json.Marshal(pos)
	assert.Nil(t, err)

	checker := func(b2 []byte) error {
		assert.Equal(t, b1, b2)
		return nil
	}

	producerError := errors.New("produce message was failed because producer falled asleep")

	saramaProducer.ExpectSendMessageWithCheckerFunctionAndFail(checker, producerError)

	err = consumerGroupHandler.controller.handleMessage(pos)
	assert.EqualError(t, err, "produce message was failed because producer falled asleep")
}

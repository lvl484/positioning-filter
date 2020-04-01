// Package kafka provides producer and consumer to work with kafka topics
package kafka

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	saramaMocks "github.com/Shopify/sarama/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	mockKafka "github.com/lvl484/positioning-filter/kafka/mocks"
	mockMatcher "github.com/lvl484/positioning-filter/matcher/mocks"
	"github.com/lvl484/positioning-filter/position"
)

func TestNewConsumerGroupHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	producer := mockKafka.NewMockProducer(ctrl)
	matcher := mockMatcher.NewMockMatcher(ctrl)

	consumerGroupHandler := newConsumerGroupHandler(matcher, producer)
	assert.NotNil(t, consumerGroupHandler.controller.matcher)
	assert.NotNil(t, consumerGroupHandler.controller.producer)
}

func TestMessageControllerHandleMessageSuccess(t *testing.T) {
	matcherCtrl := gomock.NewController(t)
	defer matcherCtrl.Finish()
	matcher := mockMatcher.NewMockMatcher(matcherCtrl)

	saramaProducer := saramaMocks.NewSyncProducer(t, nil)
	producer := &producer{
		KafkaProducer: saramaProducer,
		Config: &Config{
			ProducerTopic: "filtered-positions",
		},
	}

	consumerGroupHandler := newConsumerGroupHandler(matcher, producer)

	pos := position.Position{
		UserID:    uuid.New(),
		Latitude:  float32(1.111),
		Longitude: float32(2.222),
		Timestamp: time.Now(),
		Arrival:   time.Now().Add(time.Second),
	}

	b1, err := json.Marshal(pos)
	assert.NoError(t, err)

	checker := func(b2 []byte) error {
		assert.Equal(t, b1, b2)
		return nil
	}

	saramaProducer.ExpectSendMessageWithCheckerFunctionAndSucceed(checker)

	matcher.EXPECT().Match(pos).Return(true, nil)

	err = consumerGroupHandler.controller.handleMessage(pos)
	assert.NoError(t, err)
}

func TestMessageControllerHandleMessageFail(t *testing.T) {
	matcherCtrl := gomock.NewController(t)
	defer matcherCtrl.Finish()
	matcher := mockMatcher.NewMockMatcher(matcherCtrl)

	saramaProducer := saramaMocks.NewSyncProducer(t, nil)
	producer := &producer{
		KafkaProducer: saramaProducer,
		Config: &Config{
			ProducerTopic: "filtered-positions",
		},
	}

	consumerGroupHandler := newConsumerGroupHandler(matcher, producer)

	pos := position.Position{
		UserID:    uuid.New(),
		Latitude:  float32(1.666),
		Longitude: float32(1.333),
		Timestamp: time.Now(),
		Arrival:   time.Now().Add(time.Second),
	}

	b1, err := json.Marshal(pos)
	assert.NoError(t, err)

	checker := func(b2 []byte) error {
		assert.Equal(t, b1, b2)
		return nil
	}

	producerError := errors.New("produce message was failed because producer falled asleep")

	saramaProducer.ExpectSendMessageWithCheckerFunctionAndFail(checker, producerError)

	matcher.EXPECT().Match(pos).Return(true, nil)

	err = consumerGroupHandler.controller.handleMessage(pos)
	assert.EqualError(t, err, "produce message was failed because producer falled asleep")
}

// Package kafka provides producer and consumer to work with kafka topics
package kafka

import (
	"errors"
	"testing"
	"time"

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

	consumerGroupHandler := newConsumerGroupHandler(matcher, producer, logger)
	assert.NotNil(t, consumerGroupHandler.controller.matcher)
	assert.NotNil(t, consumerGroupHandler.controller.producer)
}

func TestMessageControllerHandleMessageSuccess(t *testing.T) {
	producerCtrl, producer, matcherCtrl, matcher, consumerGroupHandler := testConsumerGroupHandlerInit(t)
	defer producerCtrl.Finish()
	defer matcherCtrl.Finish()

	pos := newPosition()

	matcher.EXPECT().Match(pos).Return(true, nil)
	producer.EXPECT().Produce(pos).Return(nil)

	err := consumerGroupHandler.controller.handleMessage(pos)
	assert.NoError(t, err)
}

func TestMessageControllerHandleMessageFailMatcher(t *testing.T) {
	producerCtrl, _, matcherCtrl, matcher, consumerGroupHandler := testConsumerGroupHandlerInit(t)
	defer producerCtrl.Finish()
	defer matcherCtrl.Finish()

	pos := newPosition()

	matcherError := errors.New("matcher doesn't want to match, matcher want to play minecraft")

	matcher.EXPECT().Match(pos).Return(false, matcherError)

	err := consumerGroupHandler.controller.handleMessage(pos)
	assert.EqualError(t, err, "matcher doesn't want to match, matcher want to play minecraft")
}

func TestMessageControllerHandleMessageNotMatched(t *testing.T) {
	producerCtrl, _, matcherCtrl, matcher, consumerGroupHandler := testConsumerGroupHandlerInit(t)
	defer producerCtrl.Finish()
	defer matcherCtrl.Finish()

	pos := newPosition()

	matcher.EXPECT().Match(pos).Return(false, nil)

	err := consumerGroupHandler.controller.handleMessage(pos)
	assert.NoError(t, err)
}

func TestMessageControllerHandleMessageFailProducer(t *testing.T) {
	producerCtrl, producer, matcherCtrl, matcher, consumerGroupHandler := testConsumerGroupHandlerInit(t)
	defer producerCtrl.Finish()
	defer matcherCtrl.Finish()

	pos := newPosition()

	producerError := errors.New("produce message was failed because producer falled asleep")

	matcher.EXPECT().Match(pos).Return(true, nil)
	producer.EXPECT().Produce(pos).Return(producerError)

	err := consumerGroupHandler.controller.handleMessage(pos)
	assert.EqualError(t, err, "produce message was failed because producer falled asleep")
}

func newPosition() position.Position {
	return position.Position{
		UserID:    uuid.New(),
		Latitude:  float32(1.111),
		Longitude: float32(2.222),
		Timestamp: time.Now(),
		Arrival:   time.Now().Add(time.Second),
	}
}

func testConsumerGroupHandlerInit(t *testing.T) (
	*gomock.Controller, *mockKafka.MockProducer, *gomock.Controller, *mockMatcher.MockMatcher, *consumerGroupHandler,
) {
	matcherCtrl := gomock.NewController(t)
	matcher := mockMatcher.NewMockMatcher(matcherCtrl)

	producerCtrl := gomock.NewController(t)
	producer := mockKafka.NewMockProducer(producerCtrl)

	consumerGroupHandler := newConsumerGroupHandler(matcher, producer, logger)

	return producerCtrl, producer, matcherCtrl, matcher, consumerGroupHandler
}

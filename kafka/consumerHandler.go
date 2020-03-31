// Package kafka provides producer and consumer to work with kafka topics
package kafka

import (
	"encoding/json"
	"log"

	"github.com/Shopify/sarama"
	"github.com/lvl484/positioning-filter/matcher"
	"github.com/lvl484/positioning-filter/position"
)

type consumerGroupHandler struct {
	controller messageController
}

// Setup is used to make ConsumerGroupHandler implement sarama.ConsumerGroupHandler interface
func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

// Setup is used to make ConsumerGroupHandler implement sarama.ConsumerGroupHandler interface
func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim consume messages from kafka topic within consumer group
func (h consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var p position.Position
		if err := json.Unmarshal(msg.Value, &p); err != nil {
			log.Println()
			continue
		}

		if err := h.controller.handleMessage(p); err != nil {
			log.Println(err)
			continue
		}

		sess.MarkMessage(msg, "")
	}

	return nil
}

type messageController struct {
	matcher  matcher.Matcher
	producer Producer
}

func (m *messageController) handleMessage(pos position.Position) error {
	matched, err := m.matcher.Match(pos)
	if err != nil {
		return err
	}

	if matched {
		if err := m.producer.Produce(pos); err != nil {
			return err
		}
	}

	return nil
}

func newConsumerGroupHandler(matcher matcher.Matcher, producer Producer) *consumerGroupHandler {
	return &consumerGroupHandler{
		controller: messageController{
			matcher:  matcher,
			producer: producer,
		},
	}
}

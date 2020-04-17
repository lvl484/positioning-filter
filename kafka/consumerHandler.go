// Package kafka provides producer and consumer to work with kafka topics
package kafka

import (
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/lvl484/positioning-filter/matcher"
	"github.com/lvl484/positioning-filter/position"
	"github.com/sirupsen/logrus"
)

type consumerGroupHandler struct {
	controller messageController
	log        *logrus.Logger
}

// Setup is used to make ConsumerGroupHandler implement sarama.ConsumerGroupHandler interface
func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

// Setup is used to make ConsumerGroupHandler implement sarama.ConsumerGroupHandler interface
func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim consume messages from kafka topic within consumer group
func (h consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.log.Infof("Kafka message was consumed from topic positions partition %v", msg.Partition)

		var p position.Position

		sess.MarkMessage(msg, "")

		if err := json.Unmarshal(msg.Value, &p); err != nil {
			h.log.Errorf("Can't unmarshal kafka message value: %v", err)
			continue
		}

		if err := h.controller.handleMessage(p); err != nil {
			h.log.Errorf("Can't handle kafka message value: %v", err)
			continue
		}

		h.log.Infof("Kafka message for %v was successfully handled", p.UserID)
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

	if !matched {
		return nil
	}

	return m.producer.Produce(pos)
}

func newConsumerGroupHandler(matcher matcher.Matcher, producer Producer, log *logrus.Logger) *consumerGroupHandler {
	return &consumerGroupHandler{
		controller: messageController{
			matcher:  matcher,
			producer: producer,
		},
		log: log,
	}
}

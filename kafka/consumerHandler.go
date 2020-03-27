// Package kafka provides producer and consumer to work with kafka topics
package kafka

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/lvl484/positioning-filter/position"
)

type ConsumerGroupHandler struct {
	Msg chan position.Position
}

// Setup is used to make ConsumerGroupHandler implement sarama.ConsumerGroupHandler interface
func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

// Setup is used to make ConsumerGroupHandler implement sarama.ConsumerGroupHandler interface
func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim consume messages from kafka topic within consumer group
func (h ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("Message topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)

		var p position.Position
		if err := json.Unmarshal(msg.Value, &p); err != nil {
			log.Println(err)
		}

		h.Msg <- p

		sess.MarkMessage(msg, "")
	}

	return nil
}

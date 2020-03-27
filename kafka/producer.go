// Package kafka provides producer and consumer to work with kafka topics
package kafka

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/lvl484/positioning-filter/position"
)

type Producer struct {
	KafkaProducer sarama.SyncProducer
	Config        *Config
}

func NewProducer(config *Config) (*Producer, error) {
	addr := []string{fmt.Sprintf("%v:%v", config.Host, config.Port)}
	producer, err := sarama.NewSyncProducer(addr, nil)

	if err != nil {
		return nil, err
	}

	return &Producer{
		KafkaProducer: producer,
		Config:        config,
	}, nil
}

// Produce message with filtered position to kafka topic "filtered-positions"
func (p Producer) Produce(msgChan chan position.Position) {
	for {
		pos := <-msgChan
		encodedPos, err := json.Marshal(pos)

		if err != nil {
			log.Println(err)
		}

		msg := &sarama.ProducerMessage{
			Topic: p.Config.ProducerTopic,
			Key:   sarama.StringEncoder(pos.UserID.String()),
			Value: sarama.ByteEncoder(encodedPos),
		}

		if _, _, err := p.KafkaProducer.SendMessage(msg); err != nil {
			log.Println(err)
		}
	}
}

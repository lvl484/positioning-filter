package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
)

type Producer struct {
	KafkaProducer sarama.SyncProducer
}

func NewProducer(config *Config) (*Producer, error) {
	addr := []string{fmt.Sprintf("%v:%v", config.Host, config.Port)}
	producer, err := sarama.NewSyncProducer(addr, nil)

	if err != nil {
		return nil, err
	}

	return &Producer{
		KafkaProducer: producer,
	}, nil
}

func (p Producer) Produce(msg *sarama.ProducerMessage) {
	for {
		p.KafkaProducer.SendMessage(msg)
	}
}

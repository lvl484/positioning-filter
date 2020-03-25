package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
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

func (p Producer) Produce(userID string, data []byte) {
	msg := &sarama.ProducerMessage{
		Topic: p.Config.ProducerTopic,
		Key:   sarama.StringEncoder(userID),
		Value: sarama.ByteEncoder(data),
	}

	p.KafkaProducer.SendMessage(msg)
}

// Package kafka provides producer and consumer to work with kafka topics
package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/lvl484/positioning-filter/position"
)

type Producer interface {
	Produce(position.Position) error
	Close() error
}

type producer struct {
	KafkaProducer sarama.SyncProducer
	Config        *Config
}

func NewProducer(config *Config) (Producer, error) {
	addr := []string{fmt.Sprintf("%v:%v", config.Host, config.Port)}
	saramaProducer, err := sarama.NewSyncProducer(addr, nil)

	if err != nil {
		return nil, err
	}

	return &producer{
		KafkaProducer: saramaProducer,
		Config:        config,
	}, nil
}

// Produce message with filtered position to kafka topic "filtered-positions"
func (p *producer) Produce(pos position.Position) error {
	encodedPos, err := json.Marshal(pos)

	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: p.Config.ProducerTopic,
		Key:   sarama.StringEncoder(pos.UserID.String()),
		Value: sarama.ByteEncoder(encodedPos),
	}

	if _, _, err := p.KafkaProducer.SendMessage(msg); err != nil {
		return err
	}

	return nil
}

func (p *producer) Close() error {
	return p.KafkaProducer.Close()
}

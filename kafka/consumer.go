// Package kafka provides producer and consumer to work with kafka topics
package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/lvl484/positioning-filter/position"
)

const kafkaVersion = "2.4.1"

type Consumer struct {
	ConsumerGroup sarama.ConsumerGroup
	Config        *Config
}

func NewConsumer(config *Config) (*Consumer, error) {
	addr := []string{fmt.Sprintf("%v:%v", config.Host, config.Port)}
	version, err := sarama.ParseKafkaVersion(kafkaVersion)

	if err != nil {
		log.Println(err)
	}

	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = version

	consumerGroup, err := sarama.NewConsumerGroup(addr, config.ConsumerGroupID, saramaConfig)

	if err != nil {
		return nil, err
	}

	return &Consumer{
		ConsumerGroup: consumerGroup,
		Config:        config,
	}, nil
}

func (c Consumer) Consume(msgChan chan position.Position) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handler := ConsumerGroupHandler{
		Msg: msgChan,
	}

	for {
		if err := c.ConsumerGroup.Consume(ctx, []string{c.Config.ConsumerTopic}, handler); err != nil {
			log.Println(err)
		}
	}
}

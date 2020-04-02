// Package kafka provides producer and consumer to work with kafka topics
package kafka

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/lvl484/positioning-filter/matcher"
)

type Consumer struct {
	ConsumerGroup sarama.ConsumerGroup
	Config        *Config
	closeChan     chan bool
	once          sync.Once
}

func NewConsumer(config *Config) (*Consumer, error) {
	addr := getKafkaAddr(config)
	version, err := sarama.ParseKafkaVersion(config.Version)

	if err != nil {
		return nil, err
	}

	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = version
	saramaConfig.Consumer.Return.Errors = true

	consumerGroup, err := sarama.NewConsumerGroup(addr, config.ConsumerGroupID, saramaConfig)

	if err != nil {
		return nil, err
	}

	closeChan := make(chan bool)

	return &Consumer{
		ConsumerGroup: consumerGroup,
		Config:        config,
		closeChan:     closeChan,
		once:          sync.Once{},
	}, nil
}

func (c *Consumer) Consume(matcher matcher.Matcher, producer Producer) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handler := newConsumerGroupHandler(matcher, producer)

	for {
		select {
		case _, ok := <-c.closeChan:
			if !ok {
				return
			}
		case err := <-c.ConsumerGroup.Errors():
			log.Println(err)
		default:
			if err := c.ConsumerGroup.Consume(ctx, []string{c.Config.ConsumerTopic}, handler); err != nil {
				log.Println(err)
			}
		}
	}
}

// Close closes Consume func. Returns nil error to implement io.Closer interface
func (c *Consumer) Close() error {
	c.once.Do(func() { close(c.closeChan) })
	return nil
}

func getKafkaAddr(c *Config) []string {
	return []string{fmt.Sprintf("%v:%v", c.Host, c.Port)}
}

package kafka

import (
	"fmt"
	"log"

	"github.com/Shopify/sarama"
)

type Consumer struct {
	pc     sarama.PartitionConsumer
	master sarama.Consumer
}

func NewConsumer(config *Config) (*Consumer, error) {
	addr := []string{fmt.Sprintf("%v:%v", config.Host, config.Port)}
	master, err := sarama.NewConsumer(addr, nil)

	if err != nil {
		return nil, err
	}

	consumer, err := master.ConsumePartition(config.ConsumerTopic, config.ConsumerPartition, sarama.OffsetNewest)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		pc:     consumer,
		master: master}, nil
}

func (c Consumer) Consume() {
	for {
		select {
		case err := <-c.pc.Errors():
			log.Println(err)
		case msg := <-c.pc.Messages():
			log.Println(msg)
		}
	}
}

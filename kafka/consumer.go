package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
)

type Consumer struct {
	Pc     sarama.PartitionConsumer
	Master sarama.Consumer
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
		Pc:     consumer,
		Master: master}, nil
}

func (c Consumer) Consume(errchan chan *sarama.ConsumerError, msgchan chan *sarama.ConsumerMessage) {
	for {
		select {
		case err := <-c.Pc.Errors():
			errchan <- err
		case msg := <-c.Pc.Messages():
			msgchan <- msg
		}
	}
}

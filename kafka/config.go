// Package kafka provides producer and consumer to work with kafka topics
package kafka

type Config struct {
	Host            string
	Port            string
	Version         string
	ConsumerTopic   string
	ConsumerGroupID string
	ProducerTopic   string
}

// Package kafka provides producer and consumer to work with kafka topics
package kafka

// Config for kafka. Contains data to connect and handle messages
type Config struct {
	Host            string
	Port            string
	Version         string
	ConsumerTopic   string
	ConsumerGroupID string
	ProducerTopic   string
}

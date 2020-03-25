package kafka

type Config struct {
	Host              string
	Port              string
	ConsumerTopic     string
	ConsumerPartition int32
	ProducerTopic     string
}

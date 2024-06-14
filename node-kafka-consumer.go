package flownodes

import "github.com/IBM/sarama"

type KafkaConsumer struct {
	// Url amqp connect to kafka
	Url string `yaml:"url" env:"OSMO_KAFKA_URL"`
	// Key name where will be read messages
	Key string `yaml:"key"`

	// msg current message
	msg []byte

	consumer  sarama.Consumer
	partition sarama.PartitionConsumer
}

func (kc *KafkaConsumer) Init(db DB) (err error) {
	return
}

func (kc *KafkaConsumer) Execute(c Ctx) (out Resp, err error) {
	return nil, nil
}

func (kc *KafkaConsumer) GetPorts() (in, out []Port) {
	return in, out
}

package flownodes

import "github.com/IBM/sarama"

type KafkaProducer struct {
	// Url amqp connect to kafka
	Url string `yaml:"url" env:"OSMO_KAFKA_URL"`
	// Key name where will be read messages
	Key string `yaml:"key"`

	producer sarama.SyncProducer

	in Port
}

func (kp *KafkaProducer) Init(db DB) (err error) {
	return
}

func (kp *KafkaProducer) Execute(c Ctx) (out Resp, err error) {
	return nil, nil
}

func (kp *KafkaProducer) GetPorts() (in, out []Port) {
	return []Port{kp.in}, out
}

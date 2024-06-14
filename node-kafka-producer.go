package flownodes

import (
	"errors"

	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	// Url amqp connect to kafka
	Url string `yaml:"url" env:"OSMO_KAFKA_URL"`
	// Topic will be read from kafka
	Topic string `yaml:"topic"`
	// Key will be read messages from topic
	Key string `yaml:"key"`

	// client config for kafka
	client sarama.Client
	// producer connect
	producer sarama.SyncProducer

	// in bytes to send msg to kafka
	in Port
}

func (kp *KafkaProducer) Init(db DB) (err error) {
	kp.in = Port{Name: PortBytes, Type: PortBytes}

	// init kafka connection producer
	if err := kp.initKafka(); err != nil {
		return err
	}

	return
}

// initKafka initializes the kafka
func (kp *KafkaProducer) initKafka() (err error) {
	// init config kafka
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true

	// init kafka client
	client, err := sarama.NewClient([]string{kp.Url}, cfg)
	if err != nil {
		return err
	}
	kp.client = client

	// init kafka producer
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return err
	}
	kp.producer = producer

	return
}

// isConnected check is broker connected
func (kp *KafkaProducer) isConnected() bool {
	// check connected brokers
	brokerList := kp.client.Brokers()

	// if no brokers, return false
	if len(brokerList) == 0 {
		return false
	}

	// check connect broker
	conn, err := brokerList[0].Connected()
	if err != nil {
		return false
	}

	return conn
}

// recovery if something happened error connect
func (kp *KafkaProducer) recovery() error {
	return kp.initKafka()
}

// Close connection kafka
func (kp *KafkaProducer) Close() {
	kp.client.Close()
	kp.producer.Close()
}

// getCtxData get data from ctx
func (kp *KafkaProducer) getCtxData(c Ctx) ([]byte, error) {
	var body []byte

	if data, ok := c.GetData(PortBytes); ok {
		body, ok = data.([]byte)
		if !ok {
			return nil, ErrUnexpectedDataType
		}
	}
	if len(body) == 0 {
		return nil, errors.New("msg bytes not found")
	}

	return body, nil
}

func (kp *KafkaProducer) Execute(c Ctx) (out Resp, err error) {
	body, err := kp.getCtxData(c)
	if err != nil {
		return nil, err
	}

	msg := &sarama.ProducerMessage{
		Topic: kp.Topic,
		Key:   sarama.StringEncoder(kp.Key),
		Value: sarama.ByteEncoder(body),
	}

	_, _, err = kp.producer.SendMessage(msg)
	if err != nil {
		if !kp.isConnected() {
			if err := kp.recovery(); err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}

func (kp *KafkaProducer) GetPorts() (in, out []Port) {
	return []Port{kp.in}, out
}

package flownodes

import (
	"github.com/IBM/sarama"
)

type KafkaConsumer struct {
	// Url amqp connect to kafka: host:port
	Url   string `yaml:"url" env:"OSMO_KAFKA_URL"`
	Topic string `yaml:"topic"`
	// Key name where will be read messages
	Key string `yaml:"key"`

	// msg current message
	msg []byte

	client    sarama.Client
	consumer  sarama.Consumer
	partition sarama.PartitionConsumer

	// out bytes from queue kafka
	out Port
}

func (kc *KafkaConsumer) Init(db DB) (err error) {
	kc.out = Port{Name: PortBytes, Type: PortBytes}

	// init kafka connection consumer
	if err := kc.initKafka(); err != nil {
		return err
	}

	return
}

// initKafka initializes the kafka
func (kc *KafkaConsumer) initKafka() (err error) {
	// init config kafka
	cfg := sarama.NewConfig()
	cfg.Consumer.Return.Errors = true
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	// init kafka client
	client, err := sarama.NewClient([]string{kc.Url}, cfg)
	if err != nil {
		return err
	}
	kc.client = client

	// init kafka consumer
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return err
	}
	kc.consumer = consumer

	// init partition
	partitionConsumer, err := consumer.ConsumePartition(kc.Topic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	kc.partition = partitionConsumer

	// call function to read msg from kafka in background
	kc.readMsg()

	return
}

// prepareForReadMsg queue declare for read messages
func (kc *KafkaConsumer) prepareForReadMsg() error {
	return nil
}

// readMsg read a messages from ctx
func (kc *KafkaConsumer) readMsg() {
	go func() {
		for msg := range kc.partition.Messages() {
			if string(msg.Key) == kc.Key {
				kc.msg = msg.Value
			}
		}
	}()
}

func (kc *KafkaConsumer) isConnected() bool {
	// check connected brokers
	brokerList := kc.client.Brokers()

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
func (kc *KafkaConsumer) recovery() error {
	return kc.initKafka()
}

// Close connection kafka
func (kc *KafkaConsumer) Close() {
	kc.client.Close()
	kc.consumer.Close()
	kc.partition.Close()
}

func (kc *KafkaConsumer) Execute(c Ctx) (out Resp, err error) {
	if !kc.isConnected() {
		if err := kc.recovery(); err != nil {
			return nil, err
		}
	}

	if len(kc.msg) == 0 {
		return nil, nil
	}

	return Resp{kc.out.Name: kc.msg}, nil
}

func (kc *KafkaConsumer) GetPorts() (in, out []Port) {
	return in, []Port{kc.out}
}

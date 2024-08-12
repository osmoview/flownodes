package flownodes

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/osmoview/flowdata"

	"github.com/IBM/sarama"
)

func init() {
	RegisterNode("kafka-consumer", NewKafkaConsumer)
}

type KafkaConsumer struct {
	// Url amqp connect to kafka: host:port
	Url string `yaml:"url" env:"OSMO_KAFKA_URL"`
	// Topic will read from kafka
	Topic string `yaml:"topic" env:"OSMO_KAFKA_TOPIC"`
	// Key will read messages from topic
	Key string `yaml:"key"`
	// DataType type of data in queue: frame, objects, ect...
	DataType string `yaml:"data_type"`
	// DurCheckConnect timeout to check connection
	DurCheckConnect time.Duration

	// msg current message
	msg interface{}

	// client config for kafka
	client sarama.Client
	// consumer connect
	consumer sarama.Consumer
	// partition connect
	partition sarama.PartitionConsumer
	// errorChan chan errors from kafka
	errorChan <-chan *sarama.ConsumerError

	// out data read from kafka
	out []Port
}

func NewKafkaConsumer() Node {
	return &KafkaConsumer{}
}

func (kc *KafkaConsumer) Init(db DB) (err error) {
	// initialize ports
	kc.initPorts()

	if kc.Url == "" {
		return ErrUrlKafkaNotFound
	}

	if kc.Topic == "" {
		return ErrTopicKafkaNotFound
	}

	if kc.Key == "" {
		return ErrKeyKafkaNotFound
	}

	if kc.DataType == "" {
		return ErrDataTypeKafkaNotFound
	}

	if kc.DurCheckConnect == 0 {
		kc.DurCheckConnect = 30 * time.Second
	}

	// init kafka connection consumer
	if err := kc.initKafka(); err != nil {
		return err
	}

	// call function to check connection kafka in background
	kc.checkConnection()

	return
}

// initPorts initializes input ports
func (kc *KafkaConsumer) initPorts() {
	kc.out = append(kc.out, Port{Name: PortFrame, Type: PortFrame})
	kc.out = append(kc.out, Port{Name: PortObjects, Type: PortObjects})
	kc.out = append(kc.out, Port{Name: PortString, Type: PortString})
	kc.out = append(kc.out, Port{Name: PortNumber, Type: PortNumber})
	kc.out = append(kc.out, Port{Name: PortBytes, Type: PortBytes})
	kc.out = append(kc.out, Port{Name: PortAny, Type: PortAny})
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

	// init errors channel
	kc.errorChan = partitionConsumer.Errors()

	// call function to read msg from kafka in background
	go kc.readMsg()

	// call function to handle errors from kafka consumer in background
	go kc.handleErrors()

	return
}

// readMsg read a messages from ctx
func (kc *KafkaConsumer) readMsg() {
	for msg := range kc.partition.Messages() {
		if string(msg.Key) == kc.Key {
			if err := kc.parseMsg(string(msg.Headers[0].Value), msg.Value); err != nil {
				log.Error(err)
				continue
			}
		}
	}
}

// handleErrors handle errors from kafka consumer
func (kc *KafkaConsumer) handleErrors() {
	for err := range kc.errorChan {
		log.Error(err)
	}
}

// parseMsg parse a message bytes to need format
func (kc *KafkaConsumer) parseMsg(contentType string, body []byte) error {
	if contentType == ContentTypeJSON {
		if kc.DataType == PortFrame {
			var frame *flowdata.Frame
			if err := json.Unmarshal(body, &frame); err != nil {
				return err
			}
			kc.msg = frame
		}
		if kc.DataType == PortObjects {
			var objects []*flowdata.Object
			if err := json.Unmarshal(body, &objects); err != nil {
				return err
			}
			kc.msg = objects
		}
	}

	if contentType == ContentTypeText {
		if kc.DataType == PortString {
			kc.msg = string(body)
		}
		if kc.DataType == PortNumber {
			num, err := strconv.ParseFloat(string(body), 64)
			if err != nil {
				return err
			}
			kc.msg = num
		}
	}

	if contentType == ContentTypeBinary {
		if kc.DataType == PortBytes || kc.DataType == PortAny {
			kc.msg = body
		}
	}

	return nil
}

// isConnected check is broker connected
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

// checkConnection check connection kafka and recovery
func (kc *KafkaConsumer) checkConnection() {
	for {
		if !kc.isConnected() {
			if err := kc.initKafka(); err != nil {
				log.Error(err)
			}
		}
		time.Sleep(kc.DurCheckConnect)
	}
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

// getRespData get response data
func (kc *KafkaConsumer) getRespData() (out Resp, err error) {
	var resp Resp

	if kc.DataType == PortFrame {
		res, ok := kc.msg.(*flowdata.Frame)
		if !ok {
			return nil, ErrUnexpectedDataType
		}
		resp = Resp{PortFrame: res}
	}

	if kc.DataType == PortObjects {
		res, ok := kc.msg.([]*flowdata.Object)
		if !ok {
			return nil, ErrUnexpectedDataType
		}
		resp = Resp{PortObjects: res}
	}

	if kc.DataType == PortString {
		res, ok := kc.msg.(string)
		if !ok {
			return nil, ErrUnexpectedDataType
		}
		resp = Resp{PortString: res}
	}

	if kc.DataType == PortNumber {
		res, ok := kc.msg.(float64)
		if !ok {
			return nil, ErrUnexpectedDataType
		}
		resp = Resp{PortNumber: res}
	}

	if kc.DataType == PortBytes {
		res, ok := kc.msg.([]byte)
		if !ok {
			return nil, ErrUnexpectedDataType
		}
		resp = Resp{PortBytes: res}
	}

	if kc.DataType == PortAny {
		res, ok := kc.msg.([]byte)
		if !ok {
			return nil, ErrUnexpectedDataType
		}
		resp = Resp{PortAny: res}
	}

	return resp, nil
}

func (kc *KafkaConsumer) Execute(c Ctx) (out Resp, err error) {
	if !kc.isConnected() {
		if err := kc.recovery(); err != nil {
			return nil, err
		}
	}

	return kc.getRespData()
}

func (kc *KafkaConsumer) GetPorts() (in, out []Port) {
	return in, kc.out
}

func (kc *KafkaConsumer) GetForm() flowdata.Form {
	return flowdata.Form{}
}

func (KafkaConsumer) GetGroup() string {
	return ""
}

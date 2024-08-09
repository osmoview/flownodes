package flownodes

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/osmoview/flowdata"

	"github.com/IBM/sarama"
)

func init() {
	RegisterNode("kafka-producer", &KafkaProducer{})
}

var (
	ErrUrlKafkaNotFound      = errors.New("url kafka not found")
	ErrTopicKafkaNotFound    = errors.New("topic kafka not found")
	ErrKeyKafkaNotFound      = errors.New("key kafka not found")
	ErrDataTypeKafkaNotFound = errors.New("data type not found")
)

type KafkaProducer struct {
	// Url amqp connect to kafka
	Url string `yaml:"url" env:"OSMO_KAFKA_URL"`
	// Topic will write from kafka
	Topic string `yaml:"topic" env:"OSMO_KAFKA_TOPIC"`
	// Key will write messages to topic
	Key string `yaml:"key"`
	// DurCheckConnect timeout to check connection
	DurCheckConnect time.Duration

	// client config for kafka
	client sarama.Client
	// producer connect
	producer sarama.SyncProducer

	// in data send to kafka
	in []Port
}

func (kp *KafkaProducer) Init(db DB) (err error) {
	// initialize ports
	kp.initPorts()

	if kp.Url == "" {
		return ErrUrlKafkaNotFound
	}

	if kp.Topic == "" {
		return ErrTopicKafkaNotFound
	}

	if kp.Key == "" {
		return ErrKeyKafkaNotFound
	}

	if kp.DurCheckConnect == 0 {
		kp.DurCheckConnect = 30 * time.Second
	}

	// init kafka connection producer
	if err := kp.initKafka(); err != nil {
		return err
	}

	// call function to check connection kafka in background
	kp.checkConnection()

	return
}

// initPorts initializes input ports
func (kp *KafkaProducer) initPorts() {
	kp.in = append(kp.in, Port{Name: PortFrame, Type: PortFrame})
	kp.in = append(kp.in, Port{Name: PortObjects, Type: PortObjects})
	kp.in = append(kp.in, Port{Name: PortString, Type: PortString})
	kp.in = append(kp.in, Port{Name: PortNumber, Type: PortNumber})
	kp.in = append(kp.in, Port{Name: PortBytes, Type: PortBytes})
	kp.in = append(kp.in, Port{Name: PortAny, Type: PortAny})
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

// checkConnection check connection kafka and recovery
func (kp *KafkaProducer) checkConnection() {
	for {
		if !kp.isConnected() {
			if err := kp.initKafka(); err != nil {
				log.Error(err)
			}
		}
		time.Sleep(kp.DurCheckConnect)
	}
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

// sendMsg send message to kafka
func (kp *KafkaProducer) sendMsg(body []byte, contentType string) (partition int32, offset int64, err error) {
	msg := &sarama.ProducerMessage{
		Topic: kp.Topic,
		Key:   sarama.StringEncoder(kp.Key),
		Value: sarama.ByteEncoder(body),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("Content-Type"),
				Value: []byte(contentType),
			},
		},
	}

	return kp.producer.SendMessage(msg)
}

// getCtxData get data from ctx
func (kp *KafkaProducer) getCtxData(c Ctx) ([]byte, string, error) {
	var data interface{}
	var contentType string

	if d, ok := c.GetData(PortFrame); ok {
		res, ok := d.(*flowdata.Frame)
		if !ok {
			return nil, contentType, ErrUnexpectedDataType
		}
		data = res
		contentType = ContentTypeJSON
	}

	if d, ok := c.GetData(PortObjects); ok {
		res, ok := d.([]*flowdata.Object)
		if !ok {
			return nil, contentType, ErrUnexpectedDataType
		}
		data = res
		contentType = ContentTypeJSON
	}

	if d, ok := c.GetData(PortString); ok {
		res, ok := d.(string)
		if !ok {
			return nil, contentType, ErrUnexpectedDataType
		}
		data = res
		contentType = ContentTypeText
	}

	if d, ok := c.GetData(PortNumber); ok {
		res, ok := d.(float64)
		if !ok {
			return nil, contentType, ErrUnexpectedDataType
		}
		data = res
		contentType = ContentTypeText
	}

	if d, ok := c.GetData(PortBytes); ok {
		res, ok := d.([]byte)
		if !ok {
			return nil, contentType, ErrUnexpectedDataType
		}
		data = res
		contentType = ContentTypeBinary
	}

	if d, ok := c.GetData(PortAny); ok {
		res, ok := d.([]byte)
		if !ok {
			return nil, contentType, ErrUnexpectedDataType
		}
		data = res
		contentType = ContentTypeBinary
	}

	var body []byte

	if contentType == ContentTypeJSON {
		b, err := json.Marshal(data)
		if err != nil {
			return nil, contentType, err
		}
		body = b
	}

	if contentType == ContentTypeText {
		body = []byte(fmt.Sprintf("%v", data))
	}

	if contentType == ContentTypeBinary {
		if b, ok := data.([]byte); ok {
			body = b
		}
	}

	return body, contentType, nil
}

func (kp *KafkaProducer) Execute(c Ctx) (out Resp, err error) {
	body, ct, err := kp.getCtxData(c)
	if err != nil {
		return nil, err
	}

	_, _, err = kp.sendMsg(body, ct)
	if err != nil {
		if !kp.isConnected() {
			if err := kp.recovery(); err != nil {
				return nil, err
			}
			if _, _, err = kp.sendMsg(body, ct); err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}

func (kp *KafkaProducer) GetPorts() (in, out []Port) {
	return kp.in, out
}

func (kc *KafkaProducer) GetForm() flowdata.Form {
	return flowdata.Form{}
}

func (KafkaProducer) GetGroup() string {
	return ""
}

package flownodes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/osmoview/flowdata"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	ErrUrlRabbitmqNotFound      = errors.New("url amqp not found")
	ErrQueueRabbitmqNotFound    = errors.New("queue not found")
	ErrDataTypeRabbitmqNotFound = errors.New("data type not found")
)

type RabbitMQProducer struct {
	// Url amqp connect to rabbitmq
	Url string `yaml:"url" env:"OSMO_RABBITMQ_URL"`
	// Queue name where will be read messages
	Queue string `yaml:"queue" env:"OSMO_RABBITMQ_QUEUE"`
	// DurCheckConnect timeout to check connection
	DurCheckConnect time.Duration
	// DurWaitPublish max wait timeout send message
	DurWaitPublish time.Duration

	// conn connection to rabbitmq
	conn *amqp.Connection
	ch   *amqp.Channel

	// in data send to rabbitmq
	in []Port
}

func (rmqp *RabbitMQProducer) Init(db DB) (err error) {
	// initialize ports
	rmqp.initPorts()

	if rmqp.Url == "" {
		return ErrUrlRabbitmqNotFound
	}

	if rmqp.Queue == "" {
		return ErrQueueRabbitmqNotFound
	}

	if rmqp.DurCheckConnect == 0 {
		rmqp.DurCheckConnect = 30 * time.Second
	}

	if rmqp.DurWaitPublish == 0 {
		rmqp.DurWaitPublish = 5 * time.Second
	}

	// init rabbitmq connection, channel
	if err := rmqp.initRabbitMQ(); err != nil {
		return err
	}

	// call function to check connection rabbitmq in background
	rmqp.checkConnection()

	return
}

// initPorts initializes input ports
func (rmqp *RabbitMQProducer) initPorts() {
	rmqp.in = append(rmqp.in, Port{Name: PortFrame, Type: PortFrame})
	rmqp.in = append(rmqp.in, Port{Name: PortObjects, Type: PortObjects})
	rmqp.in = append(rmqp.in, Port{Name: PortString, Type: PortString})
	rmqp.in = append(rmqp.in, Port{Name: PortNumber, Type: PortNumber})
	rmqp.in = append(rmqp.in, Port{Name: PortBytes, Type: PortBytes})
	rmqp.in = append(rmqp.in, Port{Name: PortAny, Type: PortAny})
}

// initRabbitMQ initializes the rabbit mq
func (rmqp *RabbitMQProducer) initRabbitMQ() error {
	conn, err := amqp.Dial(rmqp.Url)
	if err != nil {
		return err
	}
	rmqp.conn = conn

	channel, err := rmqp.conn.Channel()
	if err != nil {
		return err
	}
	rmqp.ch = channel

	// queue declare for write messages
	_, err = rmqp.ch.QueueDeclare(
		rmqp.Queue, // queue name
		true,       // durable
		false,      // auto-delete message
		false,      // exclusive queue
		false,      // noWait
		nil,        // args
	)
	if err != nil {
		return err
	}

	return nil
}

// checkConnection check connection rabbitmq and recovery
func (rmqp *RabbitMQProducer) checkConnection() {
	for {
		if rmqp.conn.IsClosed() || rmqp.ch.IsClosed() {
			if err := rmqp.initRabbitMQ(); err != nil {
				log.Println(err)
			}
		}
		time.Sleep(rmqp.DurCheckConnect)
	}
}

// recovery if something happened error connect
func (rmqp *RabbitMQProducer) recovery() error {
	return rmqp.initRabbitMQ()
}

// Close connection rabbitmq
func (rmqp *RabbitMQProducer) Close() {
	rmqp.conn.Close()
	rmqp.ch.Close()
}

// publish send message to rabbitmq
func (rmqp *RabbitMQProducer) publish(body []byte, contentType string) error {
	ctx, cancel := context.WithTimeout(context.Background(), rmqp.DurWaitPublish)
	defer cancel()

	return rmqp.ch.PublishWithContext(
		ctx,        // context
		"",         // exchange
		rmqp.Queue, // queue
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{ // msg (publishing)
			DeliveryMode: amqp.Persistent,
			ContentType:  contentType,
			Body:         body,
		})
}

// getCtxData get data from ctx
func (rmqp *RabbitMQProducer) getCtxData(c Ctx) ([]byte, string, error) {
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

func (rmqp *RabbitMQProducer) Execute(c Ctx) (out Resp, err error) {
	body, ct, err := rmqp.getCtxData(c)
	if err != nil {
		return nil, err
	}

	err = rmqp.publish(body, ct)
	if err != nil {
		if rmqp.conn.IsClosed() || rmqp.ch.IsClosed() {
			if err := rmqp.recovery(); err != nil {
				return nil, err
			}
			if err := rmqp.publish(body, ct); err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}

func (rmqp *RabbitMQProducer) GetPorts() (in, out []Port) {
	return rmqp.in, out
}

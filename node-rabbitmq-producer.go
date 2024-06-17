package flownodes

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQProducer struct {
	// Url amqp connect to rabbitmq
	Url string `yaml:"url" env:"OSMO_RABBITMQ_URL"`
	// Queue name where will be read messages
	Queue string `yaml:"OSMO_RABBITMQ_QUEUE"`

	// conn connection to rabbitmq
	conn *amqp.Connection
	ch   *amqp.Channel

	// timeoutPublishMsg max wait timeout send message
	timeoutPublishMsg time.Duration

	// in bytes to send msg to rabbitmq
	in Port
}

func (rmqp *RabbitMQProducer) Init(db DB) (err error) {
	rmqp.in = Port{Name: PortObjects, Type: PortObjects}

	if rmqp.timeoutPublishMsg == 0 {
		rmqp.timeoutPublishMsg = 5 * time.Second
	}

	// init rabbitmq connection, channel
	if err := rmqp.initRabbitMQ(); err != nil {
		return err
	}

	return
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

	// prepare to write message
	if err := rmqp.prepareForWriteMsg(); err != nil {
		return err
	}

	return nil
}

// prepareForWriteMsg queue declare for write messages
func (rmqp *RabbitMQProducer) prepareForWriteMsg() error {
	_, err := rmqp.ch.QueueDeclare(
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

// recovery if something happened error connect
func (rmqp *RabbitMQProducer) recovery() error {
	return rmqp.initRabbitMQ()
}

// Close connection rabbitmq
func (rmqp *RabbitMQProducer) Close() {
	rmqp.conn.Close()
	rmqp.ch.Close()
}

// getCtxData get data from ctx
func (rmqp *RabbitMQProducer) getCtxData(c Ctx) ([]byte, error) {
	data, ok := c.GetData(PortObjects)
	if !ok {
		return nil, ErrUnexpectedDataType
	}

	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// publish send message to rabbitmq
func (rmqp *RabbitMQProducer) publish(body []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), rmqp.timeoutPublishMsg)
	defer cancel()

	return rmqp.ch.PublishWithContext(
		ctx,        // context
		"",         // exchange
		rmqp.Queue, // queue
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{ // msg (publishing)
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})
}

func (rmqp *RabbitMQProducer) Execute(c Ctx) (out Resp, err error) {
	body, err := rmqp.getCtxData(c)
	if err != nil {
		return nil, err
	}

	err = rmqp.publish(body)
	if err != nil {
		if rmqp.conn.IsClosed() || rmqp.ch.IsClosed() {
			if err := rmqp.recovery(); err != nil {
				return nil, err
			}
			err := rmqp.publish(body)
			if err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}

func (rmqp *RabbitMQProducer) GetPorts() (in, out []Port) {
	return []Port{rmqp.in}, out
}

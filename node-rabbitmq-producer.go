package flownodes

import (
	"context"
	"errors"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQProducer struct {
	// Url amqp connect to rabbitmq
	Url string `yaml:"url" env:"OSMO_RABBITMQ_URL"`
	// Queue name where will be read messages
	Queue string `yaml:"queue"`

	// conn connection to rabbitmq
	conn *amqp.Connection
	ch   *amqp.Channel

	// in bytes to send msg to rabbitmq
	in Port
}

func (rmqp *RabbitMQProducer) Init(db DB) (err error) {
	rmqp.in = Port{Name: PortBytes, Type: PortBytes}

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
		rmqp.Queue,
		true,
		false,
		false,
		false,
		nil,
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

func (rmqp *RabbitMQProducer) Execute(c Ctx) (out Resp, err error) {
	body, err := rmqp.getCtxData(c)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = rmqp.ch.PublishWithContext(
		ctx,
		"",
		rmqp.Queue,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})
	if err != nil {
		if rmqp.conn.IsClosed() || rmqp.ch.IsClosed() {
			if err := rmqp.recovery(); err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}

func (rmqp *RabbitMQProducer) GetPorts() (in, out []Port) {
	return []Port{rmqp.in}, out
}

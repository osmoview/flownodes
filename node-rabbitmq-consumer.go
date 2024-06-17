package flownodes

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	// Url amqp connect to rabbitmq
	Url string `yaml:"url" env:"OSMO_RABBITMQ_URL"`
	// Queue name where will be read messages
	Queue string `yaml:"queue"`

	// msg current message
	msg []byte
	// conn connection to rabbitmq
	conn *amqp.Connection
	ch   *amqp.Channel
	// delivery is channel from which messages come
	delivery <-chan amqp.Delivery

	// out bytes from queue rabbitmq
	out Port
}

func (rmqc *RabbitMQConsumer) Init(db DB) (err error) {
	rmqc.out = Port{Name: PortBytes, Type: PortBytes}

	// init rabbitmq connection, channel
	if err := rmqc.initRabbitMQ(); err != nil {
		return err
	}

	return
}

// initRabbitMQ initializes the rabbit mq
func (rmqc *RabbitMQConsumer) initRabbitMQ() (err error) {
	conn, err := amqp.Dial(rmqc.Url)
	if err != nil {
		return err
	}
	rmqc.conn = conn

	channel, err := rmqc.conn.Channel()
	if err != nil {
		return err
	}
	rmqc.ch = channel

	if err := rmqc.prepareForReadMsg(); err != nil {
		return err
	}

	// call function to read msg from RabbitMQ in background
	go rmqc.readMsg()

	return
}

// prepareForReadMsg queue declare for read messages
func (rmqc *RabbitMQConsumer) prepareForReadMsg() error {
	q, err := rmqc.ch.QueueDeclare(
		rmqc.Queue, // queue name
		true,       // durable
		false,      // auto delete message
		false,      // exclusive queue
		false,      // noWait
		nil,        // args
	)
	if err != nil {
		return err
	}

	err = rmqc.ch.Qos(1, 0, false)
	if err != nil {
		return err
	}

	rmqc.delivery, err = rmqc.ch.Consume(
		q.Name, // queue name
		"",     // consumer
		false,  // autoAck
		false,  // exclusive queue
		false,  // noLocal
		false,  // noWait
		nil,    // args
	)
	if err != nil {
		return err
	}

	return nil
}

// readMsg read a messages from ctx
func (rmqc *RabbitMQConsumer) readMsg() {
	for d := range rmqc.delivery {
		rmqc.msg = d.Body
		d.Ack(false)
	}
}

// recovery if something happened error connect
func (rmqc *RabbitMQConsumer) recovery() error {
	return rmqc.initRabbitMQ()
}

// Close connection rabbitmq
func (rmqc *RabbitMQConsumer) Close() {
	rmqc.conn.Close()
	rmqc.ch.Close()
}

func (rmqc *RabbitMQConsumer) Execute(c Ctx) (out Resp, err error) {
	if rmqc.conn.IsClosed() || rmqc.ch.IsClosed() {
		if err := rmqc.recovery(); err != nil {
			return nil, err
		}
	}

	if len(rmqc.msg) == 0 {
		return nil, nil
	}
	return Resp{rmqc.out.Name: rmqc.msg}, nil
}

func (rmqc *RabbitMQConsumer) GetPorts() (in, out []Port) {
	return in, []Port{rmqc.out}
}

package flownodes

import (
	"fmt"
	"time"

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
	rmqc.out = Port{Name: "data", Type: PortBytes}

	// init rabbitmq connection, channel
	if err := rmqc.initRabbitMQ(); err != nil {
		return err
	}

	// prepare to read message
	if err := rmqc.prepareForReadMsg(); err != nil {
		return err
	}

	// goroutine to check status connection in rabbitmq
	go rmqc.checkConnection()

	// add goroutine to read msg from RabbitMQ in background
	go rmqc.readMsg()

	return
}

// initRabbitMQ initializes the rabbit mq
func (rmqc *RabbitMQConsumer) initRabbitMQ() (err error) {
	rmqc.conn, err = amqp.Dial(rmqc.Url)
	if err != nil {
		return err
	}

	rmqc.ch, err = rmqc.conn.Channel()
	if err != nil {
		return err
	}

	return
}

// checkConnection
func (rmqc *RabbitMQConsumer) checkConnection() {
	for {
		fmt.Println("check is connected")
		if rmqc.conn.IsClosed() || rmqc.ch.IsClosed() {
			fmt.Println("connection is closed!!!")
		}
		if rmqc.ch.IsClosed() {
			fmt.Println("channel is closed!!!")
		}
		time.Sleep(5 * time.Second)
	}
}

func (rmqc *RabbitMQConsumer) prepareForReadMsg() error {
	q, err := rmqc.ch.QueueDeclare(
		rmqc.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = rmqc.ch.Qos(1, 0, false)
	if err != nil {
		return err
	}

	rmqc.delivery, err = rmqc.ch.Consume(
		q.Name,
		"",
		false,
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

// readMsg
func (rmqc *RabbitMQConsumer) readMsg() {
	var forever chan struct{}

	go func() {
		for d := range rmqc.delivery {
			rmqc.msg = d.Body
			d.Ack(false)
		}
	}()

	<-forever
}

func (rmqc *RabbitMQConsumer) Close() {
	rmqc.conn.Close()
	rmqc.ch.Close()
}

func (rmqc *RabbitMQConsumer) Execute(c Ctx) (out Resp, err error) {
	if len(rmqc.msg) == 0 {
		return nil, nil
	}
	return Resp{rmqc.out.Name: rmqc.msg}, nil
}

func (rmqc *RabbitMQConsumer) GetPorts() (in, out []Port) {
	return in, []Port{rmqc.out}
}

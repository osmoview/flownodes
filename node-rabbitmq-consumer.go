package flownodes

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/osmoview/flowdata"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	// Url amqp connect to rabbitmq
	Url string `yaml:"url" env:"OSMO_RABBITMQ_URL"`
	// Queue name where will be read messages
	Queue string `yaml:"queue" env:"OSMO_RABBITMQ_QUEUE"`
	// DataType type of data in queue: frame, objects, ect...
	DataType string `yaml:"data_type"`
	// DurCheckConnect timeout to check connection
	DurCheckConnect time.Duration

	// msg current message
	msg interface{}
	// conn connection to rabbitmq
	conn *amqp.Connection
	ch   *amqp.Channel
	// delivery is channel from which messages come
	delivery <-chan amqp.Delivery

	// out bytes from queue rabbitmq
	out []Port
}

func (rmqc *RabbitMQConsumer) Init(db DB) (err error) {
	// initialize ports
	rmqc.initPorts()

	if rmqc.Url == "" {
		return ErrUrlRabbitmqNotFound
	}

	if rmqc.Queue == "" {
		return ErrQueueRabbitmqNotFound
	}

	if rmqc.DurCheckConnect == 0 {
		rmqc.DurCheckConnect = 30 * time.Second
	}

	// init rabbitmq connection, channel
	if err := rmqc.initRabbitMQ(); err != nil {
		return err
	}

	// call function to check connection rabbitmq in background
	rmqc.checkConnection()

	return
}

// initPorts initializes output ports
func (rmqc *RabbitMQConsumer) initPorts() {
	rmqc.out = append(rmqc.out, Port{Name: PortFrame, Type: PortFrame})
	rmqc.out = append(rmqc.out, Port{Name: PortObjects, Type: PortObjects})
	rmqc.out = append(rmqc.out, Port{Name: PortString, Type: PortString})
	rmqc.out = append(rmqc.out, Port{Name: PortNumber, Type: PortNumber})
	rmqc.out = append(rmqc.out, Port{Name: PortBytes, Type: PortBytes})
	rmqc.out = append(rmqc.out, Port{Name: PortAny, Type: PortAny})
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

	// call function to read msg from rabbitMQ in background
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

	err = rmqc.ch.Qos(
		1,     // prefetchCount
		0,     // prefetchSize
		false, // global
	)
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
		if err := rmqc.parseMsg(d.ContentType, d.Body); err != nil {
			log.Println(err)
			d.Ack(false)
			continue
		}
		d.Ack(false)
	}
}

// parseMsg parse a message bytes to need format
func (rmqc *RabbitMQConsumer) parseMsg(contentType string, body []byte) error {
	if contentType == ContentTypeJSON {
		if rmqc.DataType == PortFrame {
			var frame *flowdata.Frame
			if err := json.Unmarshal(body, &frame); err != nil {
				return err
			}
			rmqc.msg = frame
		}
		if rmqc.DataType == PortObjects {
			var objects []*flowdata.Object
			if err := json.Unmarshal(body, &objects); err != nil {
				return err
			}
			rmqc.msg = objects
		}
	}

	if contentType == ContentTypeText {
		if rmqc.DataType == PortString {
			rmqc.msg = string(body)
		}
		if rmqc.DataType == PortNumber {
			num, err := strconv.ParseFloat(string(body), 64)
			if err != nil {
				return err
			}
			rmqc.msg = num
		}
	}

	if contentType == ContentTypeBinary {
		if rmqc.DataType == PortBytes || rmqc.DataType == PortAny {
			rmqc.msg = body
		}
	}

	return nil
}

// checkConnection check connection rabbitmq and recovery
func (rmqc *RabbitMQConsumer) checkConnection() {
	for {
		if rmqc.conn.IsClosed() || rmqc.ch.IsClosed() {
			err := rmqc.initRabbitMQ()
			if err != nil {
				log.Println(err)
			}
		}
		time.Sleep(rmqc.DurCheckConnect)
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

// getRespData get response data
func (rmqc *RabbitMQConsumer) getRespData() (out Resp, err error) {
	var resp Resp

	if rmqc.DataType == PortFrame {
		res, ok := rmqc.msg.(*flowdata.Frame)
		if !ok {
			return nil, ErrUnexpectedDataType
		}
		resp = Resp{PortFrame: res}
	}

	if rmqc.DataType == PortObjects {
		res, ok := rmqc.msg.([]*flowdata.Object)
		if !ok {
			return nil, ErrUnexpectedDataType
		}
		resp = Resp{PortObjects: res}
	}

	if rmqc.DataType == PortString {
		res, ok := rmqc.msg.(string)
		if !ok {
			return nil, ErrUnexpectedDataType
		}
		resp = Resp{PortString: res}
	}

	if rmqc.DataType == PortNumber {
		res, ok := rmqc.msg.(float64)
		if !ok {
			return nil, ErrUnexpectedDataType
		}
		resp = Resp{PortNumber: res}
	}

	if rmqc.DataType == PortBytes {
		res, ok := rmqc.msg.([]byte)
		if !ok {
			return nil, ErrUnexpectedDataType
		}
		resp = Resp{PortBytes: res}
	}

	if rmqc.DataType == PortAny {
		res, ok := rmqc.msg.([]byte)
		if !ok {
			return nil, ErrUnexpectedDataType
		}
		resp = Resp{PortAny: res}
	}

	return resp, nil
}

func (rmqc *RabbitMQConsumer) Execute(c Ctx) (out Resp, err error) {
	if rmqc.conn.IsClosed() || rmqc.ch.IsClosed() {
		if err := rmqc.recovery(); err != nil {
			return nil, err
		}
	}

	return rmqc.getRespData()
}

func (rmqc *RabbitMQConsumer) GetPorts() (in, out []Port) {
	return in, rmqc.out
}

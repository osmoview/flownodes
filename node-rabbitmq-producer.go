package flownodes

type RabbitMQProducer struct{}

func (rmqp *RabbitMQProducer) Init(db DB) (err error) {
	return
}

func (rmqp *RabbitMQProducer) Execute(c Ctx) (out Resp, err error) {
	return nil, nil
}

func (rmqp *RabbitMQProducer) GetPorts() (in, out []Port) {
	return in, out
}

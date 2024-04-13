package flownodes

type Text struct {
	Text string `yaml:"text"`

	port Port
}

func (t *Text) Init(db DB) error {
	t.port = Port{Name: "text", Type: PortString}
	return nil
}

func (t *Text) Execute(c Ctx) (out Resp, err error) {
	return Resp{t.port.Name: t.Text}, nil
}

func (t *Text) GetPorts() (in, out []Port) {
	return nil, []Port{t.port}
}

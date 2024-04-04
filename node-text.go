package nodes

type Text struct {
	Text string

	port Port
}

func (t *Text) Init() {
	t.port = Port{Name: "text", Type: PortString}
}

func (t *Text) Execute(c Ctx) (out Resp, err error) {
	return Resp{t.port.Name: t.Text}, nil
}

func (t *Text) GetPorts() (in, out []Port) {
	return nil, []Port{t.port}
}

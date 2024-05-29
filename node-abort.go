package flownodes

const NodeAbortName = "abort"

func init() {
	RegisterNode(NodeAbortName, &Abort{})
}

type Abort struct {
	port Port
}

func (node *Abort) Init(db DB) error {
	node.port = Port{Name: "abort", Type: PortAny}
	return nil
}

func (node *Abort) Execute(c Ctx) (out Resp, err error) {
	c.Abort()
	return nil, nil
}

func (node *Abort) GetPorts() (in, out []Port) {
	return []Port{node.port}, []Port{}
}

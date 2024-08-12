package flownodes

import "github.com/osmoview/flowdata"

func init() {
	RegisterNode("abort", NewAbortNode)
}

type Abort struct {
}

func NewAbortNode() Node {
	return &Abort{}
}

func (node *Abort) Init(db DB) error {
	return nil
}

func (node *Abort) Execute(c Ctx) (out Resp, err error) {
	c.Abort()
	return nil, nil
}

func (node *Abort) GetPorts() (in, out []Port) {
	return []Port{{Name: "abort", Type: PortAny}}, []Port{}
}

func (node *Abort) GetForm() flowdata.Form {
	return flowdata.Form{}
}

func (node *Abort) GetGroup() string {
	return GroupSimple
}

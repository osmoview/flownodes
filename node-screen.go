package flownodes

import (
	"log"

	"github.com/osmoview/flowdata"
)

const NodeScreen = "screen"

func init() {
	RegisterNode(NodeScreen, &Screen{})
}

type Screen struct {
	status  Port
	frames  Port
	objects Port
}

func (s *Screen) Init(db DB) error {
	s.initPorts()

	return nil
}

func (s *Screen) initPorts() {
	s.status = Port{Name: "status", Type: PortString}
	s.frames = Port{Name: PortFrame, Type: PortFrame}
	s.objects = Port{Name: PortObjects, Type: PortObjects}
}

func (s *Screen) Execute(c Ctx) (out Resp, err error) {
	log.Println("result")
	log.Println("status:", c.MustGetData("status"))
	log.Printf("objects: %+v", c.MustGetData("objects"))

	return c.DataMap(), nil
}

func (s *Screen) GetPorts() (in, out []Port) {
	s.initPorts()
	return []Port{s.status, s.frames, s.objects}, out
}

func (s *Screen) GetForm() flowdata.Form {
	return flowdata.Form{}
}

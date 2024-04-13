package flownodes

import "log"

type Screen struct {
	status  Port
	frames  Port
	objects Port
}

func (s *Screen) Init(db DB) error {
	s.status = Port{Name: "status", Type: PortString}
	s.frames = Port{Name: PortFrame, Type: PortFrame}
	s.objects = Port{Name: PortObjects, Type: PortObjects}

	return nil
}

func (s *Screen) Execute(c Ctx) (out Resp, err error) {
	log.Println("result")
	log.Println("status:", c.MustGetData("status"))
	log.Printf("objects: %+v", c.MustGetData("objects"))

	return nil, nil
}

func (s *Screen) GetPorts() (in, out []Port) {
	return []Port{s.status, s.frames, s.objects}, out
}

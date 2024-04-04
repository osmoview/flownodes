package flownodes

import "log"

type Screen struct {
	status  Port
	frames  Port
	objects Port
}

func (s *Screen) Init() {
	s.status = Port{Name: "status", Type: PortString}
	s.frames = Port{Name: "frames", Type: PortFrames}
	s.objects = Port{Name: "objects", Type: PortObjects}
}

func (s *Screen) Execute(c Ctx) (out Resp, err error) {
	status := c.MustGetData("status")
	log.Printf("%v", status)

	return nil, nil
}

func (s *Screen) GetPorts() (in, out []Port) {
	return []Port{s.status, s.frames, s.objects}, out
}

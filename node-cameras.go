package flownodes

type Cameras struct {
	Frames Port
}

func (s *Cameras) Init() {
	s.Frames = Port{Name: "frames", Type: PortFrames}
}

func (s *Cameras) Execute(c Ctx) (out Resp, err error) {
	frame, _ := c.GetData("frames")
	// TODO: do something with data

	return Resp{"frames": frame}, nil
}

func (s *Cameras) GetPorts() (in, out []Port) {
	return in, []Port{s.Frames}
}

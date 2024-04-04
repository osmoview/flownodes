package nodes

import "github.com/osmoview/flowdata"

type Detector struct {
	ID string

	Frames  Port
	Objects Port
}

func (d *Detector) Init() {
	d.Frames = Port{Name: "frames", Type: PortFrames}
	d.Objects = Port{Name: "objects", Type: PortObjects}
}

func (d *Detector) Execute(c Ctx) (out Resp, err error) {
	// TODO: send frame to the remote service and receive some objects
	data := []flowdata.Object{{ID: "123", Type: "test"}}

	return Resp{"objects": data}, nil
}

func (d *Detector) GetPorts() (in, out []Port) {
	return []Port{d.Frames}, []Port{d.Objects}
}

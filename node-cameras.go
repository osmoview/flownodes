package flownodes

import (
	"github.com/osmoview/flowdata"
)

type Camera struct {
	ID  string `yaml:"id"`
	cam FrameReceiver

	frames Port
}

// Init cameras node, lookup camera by specified id from the `DB`
func (node *Camera) Init(db DB) (err error) {
	node.frames = Port{Name: PortFrame, Type: PortFrame}

	if node.ID != "" {
		node.cam, err = db.GetCamera(node.ID)
	}

	return
}

func (node *Camera) Execute(c Ctx) (out Resp, err error) {
	// check `entry` data as `flowdata.Frame`, if it exists, then return it
	if data, ok := c.GetData(PortFrame); ok {
		frame, ok := data.(*flowdata.Frame)
		if !ok {
			return nil, ErrUnexpectedDataType
		}

		return Resp{node.frames.Name: frame}, nil
	}

	// getting a frame from the required camera if it exists
	if node.cam != nil {
		frame, err := node.cam.GetFrame()

		// TODO: set source info from the FrameReceiver
		// TODO: return source info as port

		return Resp{node.frames.Name: frame}, err
	}

	return Resp{node.frames.Name: nil}, ErrNoDataSource
}

func (node *Camera) GetPorts() (in, out []Port) {
	return in, []Port{node.frames}
}

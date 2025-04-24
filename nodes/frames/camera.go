package nodes

import (
	"github.com/osmoview/flowdata"
	"github.com/osmoview/flownodes"
)

// Camera node, can work in the modes, as entry node, when frames incoming from the camera
// or how standalone node, and call GetFrame of the required camera
type Camera struct {
	ID string

	ds flownodes.DataSource
}

func (cam *Camera) Init(s flownodes.Storage) (err error) {
	if cam.ID != "" {
		cam.ds, err = s.GetDataSource(cam.ID)
	}

	return
}

func (cam *Camera) Run(data CameraData) (CameraData, error) {
	if cam.ds != nil {
		frame, err := cam.ds.GetFrame()
		return CameraData{Frame: frame}, err
	}

	return data, nil
}

type CameraData struct {
	Frame *flowdata.Frame
}

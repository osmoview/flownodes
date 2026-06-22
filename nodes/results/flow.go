package results

import (
	"errors"
	"time"

	"github.com/osmoview/flownodes"
)

type Flow struct {
	ID string

	flow flownodes.Flow
}

func (node *Flow) Init(storage flownodes.Storage) (err error) {
	node.flow, err = storage.GetFlow(node.ID)
	return
}

func (node *Flow) Run(data flownodes.Data, c flownodes.Ctx) (flownodes.Data, error) {
	if node.flow == nil {
		return nil, errors.New("no flow")
	}

	if data == nil {
		data = flownodes.Data{}
	}
	// substitute datasource info to the entry data
	if _, ok := data[flownodes.PortSourceInfo]; !ok {
		data[flownodes.PortSourceInfo] = c.GetSourceInfo()
	}

	ch := make(chan flownodes.NodeRunInfo, 32)
	fr, err := node.flow.Execute(data, flownodes.FlowExecuteOptions{Force: true}, ch)
	if err != nil {
		return nil, err
	}

	// resolve execution timeout from the flow txtimeout parameter
	t := time.Now().Add(node.flow.WaitTimeout(5 * time.Second))

	for {
		select {
		case result := <-ch:
			return result.OutputData, nil
		case <-time.Tick(time.Second):
			if fr.Error != nil {
				return nil, fr.Error
			}
			if time.Now().After(t) {
				return nil, errors.New("timed out")
			}
		}
	}
}

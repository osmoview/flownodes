package results

import (
	"errors"
	"fmt"
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

	ch := make(chan flownodes.NodeRunInfo, 64)
	state, err := node.flow.Execute(data, flownodes.FlowExecuteOptions{Force: true}, ch)
	if err != nil {
		return nil, err
	}
	if state.Error != nil {
		return nil, state.Error
	}

	// resolve execution timeout from the flow txtimeout parameter
	t := time.Now().Add(node.flow.WaitTimeout(10 * time.Second))

	for {
		select {
		case result := <-ch:
			if err = result.GetError(); err != nil {
				return result.OutputData, err
			}

			if result.Completed {
				return result.OutputData, nil
			}

			c.Return(result.OutputData)

		case <-time.Tick(time.Second):
			if time.Now().After(t) {
				return nil, errors.New("timed out")
			}
			fr, ok := node.flow.GetFlowRun(state.ID)
			if !ok {
				return nil, fmt.Errorf("flowrun %s not found", state.ID)
			}
			if fr.IsCompleted() {
				return nil, nil
			}
		}
	}
}

package results

import (
	"errors"
	"log"
	"time"

	"github.com/osmoview/flownodes"
)

type Flow struct {
	ID string

	flow flownodes.Flow
}

// Init
func (node *Flow) Init(storage flownodes.Storage) (err error) {
	node.flow, err = storage.GetFlow(node.ID)
	return
}

func (node *Flow) Run(data flownodes.Data) (flownodes.Data, error) {
	ch := make(chan flownodes.NodeRunInfo)
	fr, err := node.flow.Execute(data, flownodes.FlowExecuteOptions{Force: true}, ch)
	if err != nil {
		return nil, err
	}

	// TODO: configurable timeout
	t := time.Now().Add(30 * time.Second)

	for {
		select {
		case result := <-ch:
			log.Printf("%+v", result.OutputData)
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

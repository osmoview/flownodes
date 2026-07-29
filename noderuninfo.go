package flownodes

import (
	"errors"
	"time"
)

type NodeRunInfo struct {
	// flow run ID, usualy should be same as incoming frame
	ID string `json:"id"`

	// node ID, for which this run info was obtained
	NodeID string `json:"node_id"`

	// time when the node was executed
	Time time.Time `json:"time,omitempty"`

	// time when flow run was started
	FlowTime time.Time `json:"flow_time,omitempty"`

	// execution latency in ms
	Latency int `json:"latency"`

	// incoming data to the node
	InputData Data `json:"inputdata,omitempty"`

	// response from the node
	OutputData Data `json:"outputdata,omitempty"`

	// error description
	Error string `json:"error,omitempty"`

	// completed flag is true when flow is totally completed
	Completed bool `json:"completed,omitempty"`
}

func (info *NodeRunInfo) GetError() error {
	if info.Error == "" {
		return nil
	}

	return errors.New(info.Error)
}

package database

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/osmoview/flowdata"
	"github.com/osmoview/flownodes"
)

type SaveEvent struct {
	// timeout in seconds between save events
	Timeout  int `yaml:"timeout"`
	saveTime time.Time

	// alert flag
	Alert bool `yaml:"alert"`

	Group string `yaml:"group"`

	Level int `yaml:"level"`

	db flownodes.Database
}

type SaveEventData struct {
	Frame   *flowdata.Frame
	Objects []*flowdata.Object
	Regions []*flowdata.Region
	Lines   []*flowdata.Line
	Status  string
	State   string
	Group   string
}

func (node *SaveEvent) Init(db flownodes.Database) {
	node.db = db
}

func (node *SaveEvent) Run(c flownodes.Ctx, data SaveEventData) (flownodes.Data, error) {
	t := time.Now()

	if node.saveTime.After(t) {
		return nil, nil
	}

	node.saveTime = t.Add(time.Duration(node.Timeout) * time.Second)

	e := &flowdata.Event{
		ID:         uuid.NewString(),
		DataSource: c.GetSourceInfo().ID,
		Flow:       c.GetFlowInfo().ID,
		Regions:    data.Regions,
		Lines:      data.Lines,
		Status:     data.Status,
		State:      data.State,
		Objects:    data.Objects,
		Alert:      node.Alert,
		Group:      data.Group,
		Level:      node.Level,
	}

	if err := e.Save(data.Frame); err != nil {
		return nil, fmt.Errorf("failed to save event data on disk: %v", err)
	}

	return flownodes.Data{"event": e}, node.db.Add(e)
}

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

	AllowEmpty bool `yaml:"allowempty"`
	EventPerObject bool `yaml:"oneeventperobject"`

	Group string `yaml:"group"`

	db flownodes.Database
}

type SaveEventData struct {
	Frame   *flowdata.Frame
	Objects []*flowdata.Object
	Regions []*flowdata.Region
	Lines   []*flowdata.Line
	Status  string
	State   string
}

func (node *SaveEvent) Init(db flownodes.Database) {
	node.db = db
}

func (node *SaveEvent) Run(c flownodes.Ctx, data SaveEventData) (flownodes.Data, error) {
	t := time.Now()

	if node.saveTime.After(t) {
		return nil, nil
	}

	// Skip event creation if empty objects are not allowed and objects list is empty
	if !node.AllowEmpty && len(data.Objects) == 0 {
		return flownodes.Data{}, nil
	}

	node.saveTime = t.Add(time.Duration(node.Timeout) * time.Second)

	if node.EventPerObject && len(data.Objects) > 0 {
		events := make([]*flowdata.Event, 0, len(data.Objects))
		for _, obj := range data.Objects {
			e := &flowdata.Event{
				ID:         uuid.NewString(),
				DataSource: c.GetSourceInfo().ID,
				Flow:       c.GetFlowInfo().ID,
				Regions:    data.Regions,
				Lines:      data.Lines,
				Status:     data.Status,
				State:      data.State,
				Objects:    []*flowdata.Object{obj}, // single object
				Alert:      node.Alert,
				Group:      node.Group,
			}

			if err := e.Save(data.Frame); err != nil {
				return nil, fmt.Errorf("failed to save event data on disk: %v", err)
			}

			if err := node.db.Add(e); err != nil {
				return nil, err
			}

			events = append(events, e)
		}

		return flownodes.Data{"events": events}, nil
	}

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
		Group:      node.Group,
	}

	if err := e.Save(data.Frame); err != nil {
		return nil, fmt.Errorf("failed to save event data on disk: %v", err)
	}

	return flownodes.Data{"event": e}, node.db.Add(e)
}

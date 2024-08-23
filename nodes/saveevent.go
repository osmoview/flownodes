package nodes

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"github.com/osmoview/flowdata"
	"github.com/osmoview/flownodes"
)

var (
	JPEGExt          = ".jpg"
	DefaultEventsDir = "data/events" // default dir
)

type SaveEvent struct {
	// Dir specified the path to the directory to save images,
	// if empty, used default dir
	Dir string `yaml:"dir" env:"OSMO_DATABASE_DIR"`
}

type SaveEventData struct {
	Frame   *flowdata.Frame
	Objects []*flowdata.Object
	Status  string
}

func (node *SaveEvent) Init() {
	if dir := os.Getenv("OSMO_DATABASE_DIR"); dir != "" {
		node.Dir = dir
	}

	if node.Dir == "" {
		node.Dir = DefaultEventsDir
	}
}

func (node *SaveEvent) Run(c flownodes.Ctx, data SaveEventData) error {
	e := &flowdata.Event{
		ID:         uuid.NewString(),
		DataSource: c.GetSourceInfo().ID,
		Flow:       c.GetFlowInfo().ID,
		Status:     data.Status,
		Objects:    data.Objects,
	}

	if data.Frame != nil && node.Dir != "" {
		if err := node.saveFrame(data.Frame, e); err != nil {
			return fmt.Errorf("failed to save frame image: %v", err)
		}
	}

	return c.Database().Add(e)
}

func (node *SaveEvent) saveFrame(frame *flowdata.Frame, e *flowdata.Event) error {
	dir := filepath.Join(node.Dir, e.DataSource)
	fn := filepath.Join(dir, e.ID) + JPEGExt

	// nolint:errcheck
	os.MkdirAll(dir, 0755) // no need to check error, write file return it

	return os.WriteFile(fn, frame.GetBytes(), 0755)
}

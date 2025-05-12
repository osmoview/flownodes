package database

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

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

	// timeout in seconds between save events
	Timeout  int `yaml:"timeout"`
	saveTime time.Time

	db flownodes.Database
}

type SaveEventData struct {
	Frame   *flowdata.Frame
	Objects []*flowdata.Object
	Status  string
}

func (node *SaveEvent) Init(db flownodes.Database) {
	if dir := os.Getenv("OSMO_DATABASE_DIR"); dir != "" {
		node.Dir = dir
	}

	if node.Dir == "" {
		node.Dir = DefaultEventsDir
	}

	if node.Timeout == 0 {
		node.Timeout = 10
	}

	node.db = db
}

func (node *SaveEvent) Run(c flownodes.Ctx, data SaveEventData) error {
	t := time.Now()

	if node.saveTime.After(t) {
		return nil
	}

	node.saveTime = t.Add(time.Duration(node.Timeout) * time.Second)

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

	return node.db.Add(e)
}

func (node *SaveEvent) saveFrame(frame *flowdata.Frame, e *flowdata.Event) error {
	imgdata, err := frame.GetBytes()
	if err != nil {
		return err
	}

	dir := filepath.Join(node.Dir, e.DataSource)
	fn := filepath.Join(dir, e.ID) + JPEGExt

	// nolint:errcheck
	os.MkdirAll(dir, 0755) // no need to check error, write file return it

	return os.WriteFile(fn, imgdata, 0755)
}

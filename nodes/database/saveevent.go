package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"github.com/osmoview/flowdata"
	"github.com/osmoview/flownodes"
)

var (
	JPEGExt          = ".jpg"
	JSONExt          = ".json"
	DefaultEventsDir = "data/events" // default dir
)

type SaveEvent struct {
	// Dir specified the path to the directory to save images,
	// if empty, used default dir
	Dir string `yaml:"dir" env:"OSMO_DATABASE_DIR"`

	// timeout in seconds between save events
	Timeout  int `yaml:"timeout"`
	saveTime time.Time

	// alert flag
	Alert bool `yaml:"alert"`

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
		Alert:      node.Alert,
	}

	abs, rel := node.buildFilename(e)

	if data.Frame != nil && node.Dir != "" {
		e.Image = rel + JPEGExt
		if err := node.saveFrame(abs+JPEGExt, data.Frame); err != nil {
			return fmt.Errorf("failed to save frame image: %v", err)
		}
	}

	if len(data.Objects) > 0 {
		if err := node.saveObjects(abs+JSONExt, data.Objects); err != nil {
			log.Printf("failed to save objects: %v", err)
		}
	}

	return node.db.Add(e)
}

func (node *SaveEvent) saveFrame(fn string, frame *flowdata.Frame) error {
	imgdata, err := frame.GetBytes()
	if err != nil {
		return err
	}

	return os.WriteFile(fn, imgdata, 0666)
}

func (node *SaveEvent) saveObjects(fn string, objects []*flowdata.Object) error {
	b, err := json.Marshal(objects)
	if err != nil {
		return err
	}

	return os.WriteFile(fn, b, 0666)
}

func (node *SaveEvent) buildFilename(e *flowdata.Event) (abs string, rel string) {
	// relative path to the directory where save files
	reldir := filepath.Join(e.DataSource, e.Flow)

	// nolint:errcheck
	os.MkdirAll(filepath.Join(node.Dir, reldir), 0755)

	// absolute and relative path to the file, without extension
	abs = filepath.Join(node.Dir, reldir, e.ID)
	rel = filepath.Join(reldir, e.ID)

	return
}

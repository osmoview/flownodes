package flownodes

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/osmoview/flowdata"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DatabaseEventsFramesExt = ".jpg"
	DatabaseEventsFramesDir = "data/events/" // default dir
)

func init() {
	RegisterNode("database", &Database{})
}

type Database struct {
	// DSN to connect to the postgres database, like:
	// host=localhost port=5432 dbname=gorm user=gorm password=gorm sslmode=disable TimeZone=Asia/Sigapore
	DSN string `yaml:"dsn" env:"OSMO_DATABASE_DSN"`

	// Dir specified the path to the directory to save images,
	// if empty, used default dir
	Dir string `yaml:"dir" env:"OSMO_DATABASE_DIR"`

	gorm *gorm.DB

	status  Port
	frames  Port
	objects Port
}

func (db *Database) Init(_ DB) (err error) {
	db.initPorts()

	if db.Dir == "" {
		db.Dir = DatabaseEventsFramesDir
	}

	return db.initDB()
}

func (db *Database) initPorts() {
	db.status = Port{Name: "status", Type: PortString}
	db.frames = Port{Name: PortFrame, Type: PortFrame}
	db.objects = Port{Name: PortObjects, Type: PortObjects}
}

func (db *Database) initDB() (err error) {
	// prepare gorm postgres config
	cfg := postgres.Config{DSN: db.DSN, PreferSimpleProtocol: true}

	// dial to the database
	db.gorm, err = gorm.Open(postgres.New(cfg))
	if err != nil {
		return err
	}

	return db.gorm.AutoMigrate(&flowdata.Event{})
}

func (db *Database) Execute(c Ctx) (Resp, error) {
	e := &flowdata.Event{
		ID:      c.GetID(),
		Source:  c.GetSourceInfo().ID,
		Status:  fmt.Sprint(c.MustGetData(db.status.Name)),
		Objects: c.MustGetData(db.objects.Name).([]*flowdata.Object),
	}

	if frame, ok := c.GetData(db.frames.Name); ok && db.Dir != "" {
		if err := db.saveFrame(frame, e); err != nil {
			log.Printf("failed to save frame image: %v", err)
		}
	}

	return Resp{}, db.gorm.Create(e).Error
}

func (db *Database) saveFrame(v any, e *flowdata.Event) error {
	frame, ok := v.(*flowdata.Frame)
	if !ok {
		return ErrUnexpectedDataType
	}

	dir := filepath.Join(db.Dir, e.Source)
	fn := filepath.Join(dir, e.ID) + DatabaseEventsFramesExt

	// nolint:errcheck
	os.MkdirAll(dir, 0755) // no need to check error, write file return it

	return os.WriteFile(fn, frame.GetBytes(), 0755)
}

func (db *Database) GetPorts() (in, out []Port) {
	db.initPorts()
	return []Port{db.frames, db.objects, db.status}, []Port{}
}

func (db *Database) GetForm() flowdata.Form {
	return flowdata.Form{
		{Name: "dsn", Help: flowdata.L10n{"en": "Database DSN"}, Level: flowdata.FieldLevelMaster, Hook: flowdata.FieldHookDatabase},
		{Name: "dir", Help: flowdata.L10n{"en": "Frames directory"}, Level: flowdata.FieldLevelMaster, Default: DatabaseEventsFramesDir},
	}
}

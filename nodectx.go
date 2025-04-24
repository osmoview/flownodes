package flownodes

import (
	"errors"
	"time"

	"github.com/osmoview/flowdata"
)

var (
	ErrNodeNotFound        = errors.New("node not found")
	ErrNoDataSource        = errors.New("no data source")
	ErrUnexpectedDataType  = errors.New("unexpected data type")
	ErrInterpreterNotExist = errors.New("interpreter not exist")
)

const (
	NodeEntry    = "entry"
	NodeScreen   = "screen"
	NodeService  = "service"
	NodeDatabase = "database"
)

const (
	SourceTypeFlow    = "flow"
	SourceTypeService = "service"
	SourceTypeCamera  = "camera"
)

const (
	GroupSimple        = "simple"
	GroupFrames        = "frames"
	GroupObjects       = "objects"
	GroupResults       = "results"
	GroupServices      = "services"
	GroupScripts       = "scripts"
	GroupDataSources   = "datasources"
	GroupNotifications = "notifications"
	GroupFiles         = "files"
)

//
// context and response
//

// Ctx interface contains the data from incoming transations to the node
// GetData() and MustGetData() returns data from the required 'port'
// DataMap() returns all data as map[string]any
type Ctx interface {
	GetID() string
	GetTime() time.Time
	GetSourceInfo() SourceInfo
	GetFlowInfo() SourceInfo
	GetData(port string) (data any, ok bool)
	MustGetData(port string) (data any)
	DataMap() map[string]any
	Abort()
}

//
// services and datasources storage
//

type Storage interface {
	GetService(id string) (RemoteService, error)
	GetCamera(id string) (FrameReceiver, error)
	GetDataSource(id string) (DataSource, error)
	GetFlow(id string) (Flow, error)
}

type RemoteService interface {
	GetName() string
	GetAddresses(path string) ([]string, error)
	GetAddress(path string) (string, error)
	Upload(path string, data map[string]any, resp any) error
	Post(path string, data, resp any) error
	Get(path string, resp any) error
}

// Deperecated: use DataSource instead
type FrameReceiver interface {
	GetFrame() (*flowdata.Frame, error)
	GetRegions(name string) []flowdata.Region
	GetLines(name string) []flowdata.Line
}

type DataSource interface {
	GetName() string
	GetAddr() string
	GetFrame() (*flowdata.Frame, error)
	GetRegions(name string) []flowdata.Region
	GetLines(name string) []flowdata.Line
}

type Flow interface {
	GetName() string
}

//
// database
//

type Database interface {
	Add(x DatabaseTableModel) error
	Set(x DatabaseTableModel) error
	Get(x DatabaseTableModel) error
	Del(x DatabaseTableModel) error
	Upd(where any, x DatabaseTableModel) error
	Range(x DatabaseTableModel, f func(v DatabaseTableModel) bool) error
	Table(x TableScheme) (DatabaseTable, error)
	// SQL(query string, x ...any) (*sql.DB, error)
}

type DatabaseTableModel interface {
	TableName() string
}

type DatabaseTable interface {
	Select(where map[string]any) (records []Data, err error)
	SelectOne(where map[string]any) (record Data, err error)
	Insert(record Data) error
	Update(where, values map[string]any) error
	Delete(where map[string]any) error
	Range(where map[string]any, f func(Data) bool) error
}

type TableScheme struct {
	Table   string
	Columns []TableColumn
}

type TableColumn struct {
	Name string
	Type string
}

//
// source info
//

// SourceInfo contins some basic information about camera or remote service
type SourceInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

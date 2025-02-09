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
}

type RemoteService interface {
	Send(data, resp any, opt ...any) error
	GetPorts() (input, output []Port)
}

type FrameReceiver interface {
	GetFrame() (*flowdata.Frame, error)
	GetRegions(name string) []flowdata.Region
	GetLines(name string) []flowdata.Line
}

//
// database
//

type Database interface {
	Add(x DatabaseTableModel) error
	Get(x DatabaseTableModel) error
	Del(x DatabaseTableModel) error
	Range(x DatabaseTableModel, f func(v DatabaseTableModel) bool) error
}

type DatabaseTableModel interface {
	TableName() string
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

package flownodes

import (
	"errors"
	"fmt"

	"github.com/osmoview/flowdata"
)

var (
	ErrNodeNotFound       = errors.New("node not found")
	ErrNoDataSource       = errors.New("no data source")
	ErrUnexpectedDataType = errors.New("unexpected data type")
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
	SetSourceInfo(SourceInfo)
	GetSourceInfo() SourceInfo
	GetFlowInfo() SourceInfo
	GetData(port string) (data any, ok bool)
	MustGetData(port string) (data any)
	DataMap() map[string]any
	Abort()

	// get storage
	Storage() Storage
	Database() Database
}

// Data is just a `map[string]any`, for example, it used as response from the node
type Data map[string]any

// GetData from the map by the key
func (m Data) GetData(key string) (v any, ok bool) {
	v, ok = m[key]
	return
}

// MustGetData same as GetData, but do not returns second boolean variable
func (m Data) MustGetData(key string) (v any) {
	v, _ = m[key]
	return
}

// BindTo marshal-unmarshal data map to the struct by the yaml,
// yaml used, because node data stored in the yaml format
func (m Data) BindTo(x any) error {
	return fmt.Errorf("not implemented")
}

//
// services and datasources storage
//

type Storage interface {
	GetService(id string) (RemoteService, error)
	GetCamera(id string) (FrameReceiver, error)
}

type RemoteService interface {
	Send(data, resp any) error
	GetPorts() (input, output []Port)
}

type FrameReceiver interface {
	GetFrame() (*flowdata.Frame, error)
}

//
// database
//

type Database interface {
	Add(x DatabaseTableModel) error
	Get(x DatabaseTableModel) error
	Del(x DatabaseTableModel) error
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

package flownodes

import (
	"errors"
	"sync"
	"time"

	"github.com/osmoview/flowdata"
)

var (
	ErrNodeNotFound        = errors.New("node not found")
	ErrNoDataSource        = errors.New("no data source")
	ErrUnexpectedDataType  = errors.New("unexpected data type")
	ErrInterpreterNotExist = errors.New("interpreter not exist")
	ErrInterrupted         = errors.New("execution interrupted")
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
	GroupStorage       = "storage"
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
	GetSourceInfo() flowdata.SourceInfo
	GetFlowInfo() flowdata.SourceInfo
	Headers() map[string]string
	GetData(port string) (data any, ok bool)
	MustGetData(port string) (data any)
	DataMap() map[string]any
	Abort()
	Publish(data any)
	Return(map[string]any)
	Cache() *sync.Map
	Logger() Logger
}

//
// services and datasources storage
//

type Storage interface {
	GetService(id string) (RemoteService, error)
	GetDataSource(id string) (DataSource, error)
	GetFlow(id string) (Flow, error)
	GetNode(id string, data Data) (Runner, error)
	Get(id string, out any) error
	Set(id string, data any) error
	Del(id string) error
}

type RemoteService interface {
	GetName() string
	GetAddresses(path string) ([]string, error)
	GetAddress(path string) (string, error)
	Upload(path string, data map[string]any, resp any) error
	Post(path string, data, resp any) error
	Get(path string, resp any) error
	Update(data map[string]any) error
}

type DataSource interface {
	GetID() string
	GetName() string
	GetAddr() string
	GetFrame() (*flowdata.Frame, error)
	GetRegions(name string) []flowdata.Region
	GetLines(name string) []flowdata.Line
	GetData() Data
	GetInfo() Data
	Update(data map[string]any) error
}

type Flow interface {
	GetID() string
	GetName() string
	Execute(data Data, opt FlowExecuteOptions, onResponse ...chan NodeRunInfo) (*FlowExecuteState, error)
	Run(data Data) (Data, error)
	GetNodeRunner(id string) (Runner, error)
	WaitTimeout(min time.Duration) time.Duration
	AbortFlowRun(frid string) error
}

type FlowExecuteOptions struct {
	// force to execute a node
	Force bool

	// skip execute node, use input data as output data
	// this is necessary to perform classical analytics when input data corresponds
	// to output data from the entry node
	Skip bool
}

type FlowExecuteState struct {
	ID    string
	Error error
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
	Range(x DatabaseTableModel, opt DatabaseRangeOpt, f func(v DatabaseTableModel) bool) error
	Table(x TableScheme) (DatabaseTable, error)
	SQL(query string, x any) error
}

type DatabaseTableModel interface {
	TableName() string
}

type DatabaseRangeOpt struct {
	Filter map[string]any
	Offset int
	Limit  int
	Order  string
}

type DatabaseTable interface {
	Select(where map[string]any) (records []Data, err error)
	SelectOne(where map[string]any) (record Data, err error)
	Insert(record Data) error
	Update(where, values map[string]any) error
	Delete(where map[string]any) error
	Range(opt DatabaseRangeOpt, f func(DatabaseTableRow) bool) error
}

type TableScheme struct {
	Table   string
	Columns []TableColumn
}

type TableColumn struct {
	Name string
	Type string
}

type DatabaseTableRow interface {
	BindTo(x any) error
}

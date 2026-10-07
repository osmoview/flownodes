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
	// returns flowrun id
	GetID() string

	// returns time when flow run was created
	GetTime() time.Time

	// basic information about the datasource
	GetSourceInfo() flowdata.SourceInfo

	// basic information about the flow
	GetFlowInfo() flowdata.SourceInfo

	// returns current node interface
	Node() Node

	// returns prepared HTTP headers
	Headers() map[string]string

	// get data by port
	GetData(port string) (data any, ok bool)
	MustGetData(port string) (data any)

	// returns all data in an untyped format
	DataMap() map[string]any

	// abort analytics execution
	Abort()

	// publish the data, it will be possible to get it through an API request.
	Publish(data any)

	// return data from the node via outgoing ports and links, bypassing the final return
	Return(map[string]any)

	// cache is a global object accessible from all nodes and all flows of this spot,
	// any data can be stored and read here
	Cache() *sync.Map
}

//
// services and datasources storage
//

type Storage interface {
	// GetService by id, returns RemoteService interface,
	// if the service not specified in the configuration file,
	// trying to get his information from the belt
	GetService(id string) (RemoteService, error)

	// GetDataSource by id, returns DataSource interface
	GetDataSource(id string) (DataSource, error)

	// GetFlow by id, returns Flow interface for flow contains in the current spot config
	GetFlow(id string) (Flow, error)

	// GetNode returns a new instance(Runner interface) of the required node, node should be exists on the dist
	GetNode(id string, data Data) (Runner, error)

	// receive data to the belt storage and unpack it to the pointer passed as second `out any` argument
	Get(id string, out any) error
	// store any data to the belt storage
	Set(id string, data any) error
	// remove data from the belt storage
	Del(id string) error
}

type RemoteService interface {
	// GetID returns remote servcie id
	GetID() string

	// GetName should returns service name
	GetName() string

	// GetAddresses returns the list of available service addresses, with a substituted path
	GetAddresses(path string) ([]string, error)

	// GetAddress returns available sevice address(sorted by balancer) with a substituted path
	GetAddress(path string) (string, error)

	// submit data to the remote service by specified path as multipath/form-data
	Upload(path string, data map[string]any, resp any) error

	// submit data to the remote servce as json
	Post(path string, data, resp any) error

	// submit GET request to the remote service
	Get(path string, resp any) error

	// Update changes any parameters of the specified service via WARD
	Update(data map[string]any) error
}

type DataSource interface {
	// GetID returns datasource id
	GetID() string

	// GetName returns datasource name
	GetName() string

	// GetAddr returns http address for remote datasource or data.addr for datasource contains in the SPOT configuration
	GetAddr() string

	// GetFrame returns current flowdata.Frame if it exists in the last data
	GetFrame() (*flowdata.Frame, error)

	// GetRegions returns regions list by specified regions name
	GetRegions(name string) []flowdata.Region

	// GetLines returns lines by specified name
	GetLines(name string) []flowdata.Line

	// custom datasource data speciried in the configuration file
	GetData() Data

	// custom information about datasource like, cpu usage, ram usage, frame format, etc
	GetInfo() Data

	// changes any parameters of the spcified datasource via WARD
	Update(data map[string]any) error
}

type Flow interface {
	// get flow id
	GetID() string

	// get current flow name
	GetName() string

	// execute flow by specfied data, the results will be transmitted throught onResponse channels
	Execute(data Data, opt FlowExecuteOptions, onResponse ...chan NodeRunInfo) (*FlowExecuteState, error)

	// execute flow with data and receive response
	Run(data Data) (Data, error)

	// GetNodeRunner returns Runner interface for required node(by nodeid) from the current flow
	GetNodeRunner(nodeid string) (Runner, error)

	// WaitTimeout returns tx_timeout or default timeout if tx_timeout not specified
	WaitTimeout(min time.Duration) time.Duration

	// interrupt the flow execution on the spcified flowrun
	AbortFlowRun(frid string) error

	// initialize new flowrun
	NewFlowRun(sourceInfo flowdata.SourceInfo) FlowRun

	// run flowrun with input data
	RunFlowRun(fr FlowRun, data map[string]any) error

	// get flowrun by flowrun id
	GetFlowRun(id string) (FlowRun, bool)

	// create and initialize clone of the flow
	Clone() (Flow, error)

	// close flow
	Close()
}

type FlowRun interface {
	// Get flowrun id
	GetID() string

	// subscribe to the flowrun results
	SubscribeOnResult(ch chan NodeRunInfo)
	// unsubscribe from the results
	UnsubscribeOnResult(ch chan NodeRunInfo)

	// interrupt flowrun
	Abort()

	// get an abort flag that the flowrun execution was interrupted
	IsAborted() bool

	// get result flag, true if the resulting node for this flowrun has been executed
	IsResulted() bool

	// IsCompleted will be true if all nodes have completed,
	// and no other nodes are waiting for execution on incoming links
	IsCompleted() bool
}

type FlowExecuteOptions struct {
	// force to execute a node
	Force bool

	// skip execute node, use input data as output data
	// this is necessary to perform classical analytics when input data corresponds
	// to output data from the entry node
	Skip bool

	// create flow clone before execute and terimante after it completed
	Clone bool
}

type FlowExecuteState struct {
	// FlowRun ID
	ID string

	// flowrun initialization error
	Error error
}

type Node interface {
	// get node id
	GetID() string

	// get input port by port name
	GetInputPort(port string) (Port, bool)
	// get output port by port name
	GetOutputPort(port string) (Port, bool)

	// get all node input and output ports
	GetPorts() (input, output []Port)
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

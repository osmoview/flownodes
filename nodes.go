package flownodes

import (
	"errors"
	"fmt"
	"strings"

	"github.com/osmoview/flowdata"
	"github.com/osmoview/logger"
	"gopkg.in/yaml.v3"
)

var (
	log = logger.New("FLOWNODES")

	// EntryNode contains default ID of the entry node, 'entry' by default
	EntryNode = "entry"

	// SourceNodes is a list of nodes that will be executed when the main entry node starts execution
	SourceNodes = []string{}

	// Nodes contains all available(registered) nodes
	// WARN: maybe required to add sync.Mutex to access to the Nodes map
	Nodes = map[string]NodeCreator{}
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

// RegisterNode add new node to the global nodes storage
func RegisterNode(node string, n NodeCreator) {
	Nodes[node] = n

	if n.GetDriver().IsDataSource {
		SourceNodes = append(SourceNodes, node)
	}
}

// LookupNode by type and returns new instance of it
func LookupNode(node string) (NodeCreator, error) {
	log.Debug("lookup node", node)

	n, ok := Nodes[node]
	if !ok {
		return nil, fmt.Errorf("%s %w", node, ErrNodeNotFound)
	}

	return n, nil
}

//
// nodes interfaces
//

// NodeEnvironmnent is an interface for different virtual machines, interpreters or other solutions to run and execute custom nodes
type NodesEnvironment interface {
	NewNodeCreator(conf string) (NodeCreator, error)
}

// NodeCreater describe interface to any entity that should creates new nodes
type NodeCreator interface {
	GetDriver() NodeDriver
	NewNode(data Data, storage Storage, db Database) (Runner, error)
}

// Runner interface required Run function to execute flow ctx with incoming transaction data
type Runner interface {
	// this function is called every time incoming data(`Ctx`) is received,
	// and should returns response(Data:map[string]any),
	// where the key is a port name, and the value contains the data
	Run(c Ctx) (out Data, err error)
}

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
	if v, ok = m[key]; ok {
		return v, ok
	}

	return m.GetCaseInsensitive(key)
}

// MustGetData same as GetData, but do not returns second boolean variable
func (m Data) MustGetData(key string) any {
	if v, ok := m[key]; ok {
		return v
	}

	v, _ := m.GetCaseInsensitive(key)
	return v
}

func (m Data) GetCaseInsensitive(key string) (v any, ok bool) {
	for k, v := range m {
		if strings.EqualFold(k, key) {
			return v, true
		}
	}

	return nil, false
}

// BindTo marshal-unmarshal data map to the struct by the yaml,
// yaml used, because node data stored in the yaml format
func (m Data) BindTo(x any) error {
	// no data no bind
	if len(m) == 0 {
		return nil
	}

	b, err := yaml.Marshal(m)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, x)
}

func (m *Data) Set(k string, v any) {
	(*m)[k] = v
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

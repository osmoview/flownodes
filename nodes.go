package flownodes

import (
	"errors"
	"fmt"

	"github.com/osmoview/flowdata"
	"github.com/osmoview/logger"
)

const (
	ContentTypeJSON   = "application/json"
	ContentTypeText   = "text/plain"
	ContentTypeBinary = "application/octet-stream"
)

const (
	SourceTypeCamera  = "camera"
	SourceTypeService = "service"
	SourceTypeFlow    = "flow"
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

var (
	ErrNodeNotFound       = errors.New("node not found")
	ErrNoDataSource       = errors.New("no data source")
	ErrUnexpectedDataType = errors.New("unexpected data type")
)

var (
	log = logger.New("FLOWNODES")

	// EntryNode contains default ID of the entry node, 'entry' by default
	EntryNode = "entry"

	// SourceNodes is a list of nodes that will be executed when the main entry node starts execution
	SourceNodes = []string{}

	// Nodes contains all available(registered) nodes
	Nodes = map[string]func() Node{}
)

type Node interface {
	// initialize node, set default values, etc
	Init(db DB) error

	// this function is called every time incoming data(`Ctx`) is received,
	// and should returns response(Resp:map[string]any),
	// where the key is a port name, and the value contains the data
	Execute(c Ctx) (out Resp, err error)

	// get input and output node ports
	GetPorts() (in, out []Port)

	GetForm() (form flowdata.Form)

	GetGroup() string
}

func RegisterNode(t string, f func() Node) {
	Nodes[t] = f
}

// LookupNode by type and returns new instance of it
func LookupNode(t string) (Node, error) {
	log.Debug("lookup node", t)

	f, ok := Nodes[t]
	if !ok {
		return nil, fmt.Errorf("%s %w", t, ErrNodeNotFound)
	}

	return f(), nil
}

//
// context and response
//

type Ctx interface {
	GetID() string
	SetSourceInfo(SourceInfo)
	GetSourceInfo() SourceInfo
	GetFlowInfo() SourceInfo
	GetData(port string) (data any, ok bool)
	MustGetData(port string) (data any)
	DataMap() map[string]any
	Abort()
}

type Resp map[string]any

func (r Resp) GetData(key string) (v any, ok bool) {
	v, ok = r[key]
	return
}

func (r Resp) MustGetData(key string) (v any) {
	v, _ = r[key]
	return
}

//
// services storage and database interface
//

type DB interface {
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

type Executor interface {
	Execute(c Ctx) (out Resp, err error)
}

//
// source info
//

// SourceInfo contins some basic information about camera or remote service
type SourceInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	// Addr string
}

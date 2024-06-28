package flownodes

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/osmoview/flowdata"
)

const (
	ContentTypeJSON   = "application/json"
	ContentTypeText   = "text/plain"
	ContentTypeBinary = "application/octet-stream"
)

var ErrNoDataSource = errors.New("no data source")
var ErrUnexpectedDataType = errors.New("unexpected data type")

// SourceNodes is a list of nodes that will be executed when the main entry node starts execution
var SourceNodes = []string{}

// Nodes contains all available(registered) nodes
var Nodes = map[string]Node{}

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
}

func RegisterNode(t string, node Node) {
	Nodes[t] = node
}

// LookupNode by type and returns new instance of it
func LookupNode(t string) (Node, error) {
	inode, ok := Nodes[t]
	if !ok {
		return nil, fmt.Errorf("node '%s' not found", t)
	}

	return reflectCopy(inode).(Node), nil
}

func reflectCopy(inode interface{}) interface{} {
	return reflect.New(reflect.TypeOf(inode).Elem()).Interface()
}

//
// context and response
//

type Ctx interface {
	GetID() string
	SetSourceInfo(SourceInfo)
	GetSourceInfo() SourceInfo
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
// ports
//

const (
	PortFrame      = "frame"
	PortObjects    = "objects"
	PortString     = "string"
	PortNumber     = "number"
	PortMap        = "map"
	PortBytes      = "bytes"
	PortAny        = "any"
	PortSourceInfo = "sourceinfo"
	PortRegion     = "region"
	PortEmpty      = ""
)

// Port on the Node, contains name and type
// should be impossible to connect ports of different types
type Port struct {
	// port name, like status, frame, objects, or something other,
	// optional, if not specified, frontend should use port type as name
	Name string `yaml:"name,omitempty" json:"name"`

	// port type defines the type of data that can be transmitted over this port
	// frontend should prohibit the connection of ports of different types to each other
	Type string `yaml:"type" json:"type"`

	// text description of the port
	Desc string `yaml:"desc,omitempty" json:"desc,omitempty"`

	// required flag
	Required bool `yaml:"required" json:"required,omitempty"`
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

// SourceInfo contins some basic information about camera or remote service
type SourceInfo struct {
	ID   string
	Name string
	Type string
	// Addr string
}

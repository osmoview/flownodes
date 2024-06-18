package flownodes

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/osmoview/flowdata"
)

var ErrNoDataSource = errors.New("no data source")
var ErrUnexpectedDataType = errors.New("unexpected data type")
var ErrUrlRabbitmqNotFound = errors.New("url amqp not found")
var ErrQueueRabbitmqNotFound = errors.New("queue not found")

const (
	ContentTypeJSON   = "application/json"
	ContentTypeText   = "text/plain"
	ContentTypeBinary = "application/octet-stream"
)

const (
	NodeScreen   = "screen"
	NodeCamera   = "camera"
	NodeService  = "service"
	NodeText     = "text"
	NodeReadFile = "readfile"
	NodeFilter   = "filter"
	NodeColor    = "color"
	NodeCrop     = "crop"
)

// EntryNodes are the nodes from which the execution of analytics can begin
var EntryNodes = []string{NodeCamera}

// SourceNodes is a list of nodes that will be executed when the main entry node starts execution
var SourceNodes = []string{NodeCamera, NodeText, NodeService}

var Nodes = map[string]Node{
	NodeScreen:   &Screen{},
	NodeCamera:   &Camera{},
	NodeService:  &Service{},
	NodeText:     &Text{},
	NodeReadFile: &ReadFile{},
	NodeFilter:   &Filter{},
	NodeColor:    &Color{},
	NodeCrop:     &Crop{},
	"cpu-temp":   &CPUtemp{},
	"expr":       &Expression{},
}

type Node interface {
	// initialize node, set default values, etc
	Init(db DB) error

	// this function is called every time incoming data(`Ctx`) is received,
	// and should returns response(Resp:map[string]any),
	// where the key is a port name, and the value contains the data
	Execute(c Ctx) (out Resp, err error)

	// get input and output node ports
	GetPorts() (in, out []Port)
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
	Name string `yaml:"name,omitempty"`

	// port type defines the type of data that can be transmitted over this port
	// frontend should prohibit the connection of ports of different types to each other
	Type string `yaml:"type"`

	// text description of the port
	Desc string `yaml:"desc,omitempty"`

	// if field flag enabled, so this is not a port, but a field for entering data
	Field bool `yaml:"field,omitempty"`

	// optional flag, if true, then port provide not necessary data
	Optional bool `yaml:"optional"`
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

package flownodes

import (
	"fmt"
	"reflect"
)

const (
	NodeScreen   = "screen"
	NodeCameras  = "cameras"
	NodeScript   = "script"
	NodeText     = "text"
	NodeFilter   = "filter"
	NodeDetector = "detector"
)

// EntryNodes are the nodes from which the execution of analytics can begin
var EntryNodes = []string{NodeCameras}

// SourceNodes is a list of nodes that will be executed when the main entry node starts execution
var SourceNodes = []string{NodeCameras, NodeText, NodeScript}

var Nodes = map[string]Node{
	NodeScreen:   &Screen{},
	NodeCameras:  &Cameras{},
	NodeScript:   &Script{},
	NodeText:     &Text{},
	NodeFilter:   &Filter{},
	NodeDetector: &Detector{},
}

type Node interface {
	// initialize node, set default values, etc
	Init()

	// this function is called every time incoming data(`Ctx``) is received,
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
	GetData(port string) (data interface{}, ok bool)
	MustGetData(port string) (data interface{})
}

type Resp map[string]interface{}

//
// ports
//

const (
	PortFrames  = "frames"
	PortObjects = "objects"
	PortString  = "string"
	PortNumber  = "number"
	PortEmpty   = ""
)

// Port on the Node, contains name and type
// should be impossible to connect ports of different types
type Port struct {
	Name string
	Type string
}

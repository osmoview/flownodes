package flownodes

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

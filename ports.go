package flownodes

const (
	// TODO: is is not all types/names of ports
	// TODO: MAYBE: separate port types and port names

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

	// PortAsterisk usually used when no requirements for specific incoming data
	PortAsterisk = "*"
)

const (
	PortModeOne    = "one"
	PortModeLot    = "lot"
	PortModeAll    = "all"
	PortModeStream = "stream"
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
	Required bool `yaml:"required,omitempty" json:"required,omitempty"`

	// port mode:
	// one - pick first data and mark port is completed, all other links to it port will be ignored
	// lot - pick at least one link, all other links will be merged, but not wait when all links executed
	// all - wait when all links executed, mode by default if mode not specified
	// stream - execute destination node directly, without mergeData
	Mode string `yaml:"mode,omitempty" json:"mode,omitempty"`
}

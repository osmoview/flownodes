package flownodes

import (
	"fmt"

	"github.com/osmoview/flowdata"
)

// NodeModule structure describe how to run and communicate with external node module
type NodeModule struct {
	ID    string        `yaml:"id" json:"id"`
	Desc  string        `yaml:"desc,omitempty" json:"desc,omitempty"`
	Label flowdata.L10n `yaml:"label,omitempty" json:"label,omitempty"`
	Node  string        `yaml:"node" json:"node"`
	Lang  string        `yaml:"lang" json:"lang"`
	Group string        `yaml:"group,omitempty" json:"group,omitempty"`

	Input  []Port `yaml:"input,omitempty" json:"input,omitempty"`
	Output []Port `yaml:"output,omitempty" json:"output,omitempty"`

	Form []flowdata.FormGroup `yaml:"form,omitempty" json:"form,omitempty"`

	// run map may contains any instruction required to the run this node
	// Run map[string]any `yaml:"run" json:"run"`

	Errors []string `yaml:"-" json:"errors"`
}

func (n NodeModule) IsValid() error {
	if n.ID == "" {
		return fmt.Errorf("no id")
	}
	if n.Node == "" {
		return fmt.Errorf("'node' is required")
	}
	if n.Lang == "" {
		return fmt.Errorf("'lang' is required")
	}

	return nil
}

// NewModuleRunner it is an interface for the some interpreter,
// virtual machine or other solution to run custom nodes.
// function `NewNode(fn string) Node` - receive path to the yaml file,
// determines how to communitcate and run this node. This function should
// returns the Node interface.
type NodeModuleRunner interface {
	NewNode(fn string) (Node, error)
}

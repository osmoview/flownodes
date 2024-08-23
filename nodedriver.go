package flownodes

import (
	"fmt"

	"github.com/osmoview/flowdata"
)

//
// INFO: NodeDriver partially overlaps with drivers for data sources and services,
//       maybe it makes sense to use a single structure
//

// NodeDriver structure describe how to run and communicate with the node
type NodeDriver struct {
	// node id, should be unique
	ID string `yaml:"id" json:"id"`

	// path to the config file
	// Path string `yaml:"-" json:"path"`

	Node string `yaml:"node" json:"node"`
	Lang string `yaml:"lang,omitempty" json:"lang,omitempty"`

	IsDataSource bool `yaml:"is_datasource,omitempty" json:"is_datasource,omitempty"`

	Desc  flowdata.L10n `yaml:"desc,omitempty" json:"desc,omitempty"`
	Label flowdata.L10n `yaml:"label,omitempty" json:"label,omitempty"`
	Group string        `yaml:"group,omitempty" json:"group,omitempty"`

	Input  []Port        `yaml:"input,omitempty" json:"input,omitempty"`
	Output []Port        `yaml:"output,omitempty" json:"output,omitempty"`
	Form   flowdata.Form `yaml:"form,omitempty" json:"form,omitempty"`

	// run map may contains any instruction required to the run this node
	// Run map[string]any `yaml:"run" json:"run"`

	Errors []string `yaml:"-" json:"errors"`
}

func (n NodeDriver) IsValid() error {
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

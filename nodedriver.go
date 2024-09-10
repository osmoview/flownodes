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
	Path string `yaml:"-" json:"path"`

	// INFO: node field is deprecated, use ID instead
	// Node string `yaml:"node" json:"-"`
	Lang string `yaml:"lang,omitempty" json:"lang,omitempty"`

	IsDataSource bool `yaml:"is_datasource,omitempty" json:"is_datasource,omitempty"`

	Desc  flowdata.L10n `yaml:"desc,omitempty" json:"desc,omitempty"`
	Label flowdata.L10n `yaml:"label,omitempty" json:"label,omitempty"`
	Group string        `yaml:"group,omitempty" json:"group,omitempty"`

	Input  []Port        `yaml:"input,omitempty" json:"input,omitempty"`
	Output []Port        `yaml:"output,omitempty" json:"output,omitempty"`
	Form   flowdata.Form `yaml:"form,omitempty" json:"form,omitempty"`

	// INFO: Run section is deprecated, each interpreter should parse it itself
	// run map may contains any instruction required to the run this node
	// Run map[string]any `yaml:"run" json:"run"`
	// Run NodeDriverRun `yaml:"run" json:"run"`

	Errors []string `yaml:"-" json:"errors"`
}

type NodeDriverRun struct {
	File   string `yaml:"file" json:"file"`
	Target string `yaml:"target" json:"target"`
	Init   string `yaml:"init,omitempty" json:"init"`
	Run    string `yaml:"run" json:"run"`
}

func (n NodeDriver) IsValid() error {
	if n.ID == "" {
		return fmt.Errorf("'id' is required and should be unique")
	}
	if n.Lang == "" {
		return fmt.Errorf("'lang' is required")
	}

	return nil
}

// func (n NodeDriver) BindRunTo(x any) error {
// 	yaml.Marshal(n.Run)
// }

package nodes

import (
	"github.com/osmoview/flowdata"
)

type ReadFile struct {
	Path string
}

type ReadFileIn struct {
	Path string
}

type ReadFileOut struct {
	File *flowdata.File `json:"file,omitempty"`
}

func (node *ReadFile) Run(in ReadFileIn) (resp ReadFileOut, err error) {
	if in.Path != "" {
		resp.File, err = flowdata.ReadFile(in.Path)
		return
	}

	if node.Path != "" {
		resp.File, err = flowdata.ReadFile(node.Path)
		return
	}

	return
}

package nodes

import (
	"fmt"

	"github.com/osmoview/flowdata"
)

type Screen struct{}

type ScreenData struct {
	Frame   *flowdata.Frame    `json:"frame"`
	Objects []*flowdata.Object `json:"objects"`
	Status  string             `json:"status"`
}

func (node *Screen) Run(data ScreenData) ScreenData {
	fmt.Printf("RESULT: %+v\n", data)
	return data
}

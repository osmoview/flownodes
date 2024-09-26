package nodes

import (
	"log"

	"github.com/osmoview/flowdata"
)

type Screen struct{}

type ScreenData struct {
	Frame   *flowdata.Frame    `json:"frame"`
	Objects []*flowdata.Object `json:"objects"`
	Status  string             `json:"status"`
}

func (node *Screen) Run(data ScreenData) ScreenData {
	log.Printf("%+v", data)
	return data
}

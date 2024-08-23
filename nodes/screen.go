package nodes

import (
	"log"

	"github.com/osmoview/flowdata"
)

type Screen struct{}

type ScreenData struct {
	Frame   *flowdata.Frame
	Objects []*flowdata.Object
	Status  string
}

func (node *Screen) Run(data ScreenData) ScreenData {
	log.Printf("%+v", data)
	return data
}

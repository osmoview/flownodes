package charts

import (
	"time"

	"github.com/osmoview/flowdata"
	"github.com/osmoview/flownodes"
)

type ChartBuildSeries struct {
	storage flownodes.Storage
	db      flownodes.Database

	Name string
	Type string
}

type ChartBuildSeriesInput struct {
	Events []flowdata.Event
	XAxis  []time.Time
}

type ChartBuildSeriesOutput struct {
	Series []ChartSeries
}

func (node *ChartBuildSeries) Init(storage flownodes.Storage, db flownodes.Database) {
	node.storage = storage
	node.db = db

	if node.Type == "" {
		node.Type = "line"
	}
}

func (node *ChartBuildSeries) Run(in ChartBuildSeriesInput) ChartBuildSeriesOutput {
	chart := ChartSeries{
		Name: node.Name,
		Type: node.Type,
		Data: make([]float64, len(in.XAxis)),
	}

	for _, e := range in.Events {
		node.setData(e, chart.Data, in.XAxis)
	}

	return ChartBuildSeriesOutput{Series: []ChartSeries{chart}}
}

func (node *ChartBuildSeries) setData(e flowdata.Event, data []float64, xaxis []time.Time) {
	for i, t := range xaxis {
		if (e.CreatedAt.After(t) || e.CreatedAt.Equal(t)) && e.CreatedAt.Before(t) {
			data[i] += float64(len(e.Objects))
			return
		}
	}
}

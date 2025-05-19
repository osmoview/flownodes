package charts

import (
	"time"

	"github.com/osmoview/echarts"
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
	Series []echarts.Series
}

func (node *ChartBuildSeries) Init(storage flownodes.Storage, db flownodes.Database) {
	node.storage = storage
	node.db = db

	if node.Type == "" {
		node.Type = "line"
	}
}

func (node *ChartBuildSeries) Run(in ChartBuildSeriesInput) ChartBuildSeriesOutput {
	chart := echarts.Series{
		Name: node.Name,
		Type: node.Type,
		Data: make([]float64, len(in.XAxis)),
	}

	for _, e := range in.Events {
		node.setData(e, chart.Data, in.XAxis)
	}

	return ChartBuildSeriesOutput{Series: []echarts.Series{chart}}
}

func (node *ChartBuildSeries) setData(e flowdata.Event, data []float64, xaxis []time.Time) {
	if len(xaxis) <= 1 {
		data[0] += float64(len(e.Objects))
		return
	}

	diff := xaxis[1].Sub(xaxis[0])

	for i, t := range xaxis {
		if (e.CreatedAt.After(t) || e.CreatedAt.Equal(t)) && e.CreatedAt.Before(t.Add(diff)) {
			data[i] += float64(len(e.Objects))
			return
		}
	}

	return
}

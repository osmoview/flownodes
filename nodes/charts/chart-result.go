package charts

import (
	"time"

	"github.com/osmoview/echarts"
)

type ChartResult struct{}

type ChartResultInput struct {
	Series []echarts.Series `json:"series"`
	XAxis  []time.Time      `json:"xAxis"`
	YAxis  echarts.Axis     `json:"yAxis"`
}

func (t *ChartResult) Run(input ChartResultInput) echarts.Echarts {
	opt := echarts.NewEcharts()
	opt.Series = input.Series
	opt.Tooltip = echarts.Tooltip{Trigger: "axis"}
	// opt.Grid = []echats.Grid
	opt.Legend = echarts.Legend{Show: echarts.TRUE}

	xaxis := echarts.Axis{}
	for _, t := range input.XAxis {
		xaxis.Data = append(xaxis.Data, t.Format(echarts.TimeLayoutDateTime))
	}

	opt.XAxis = []echarts.Axis{xaxis}
	opt.YAxis = []echarts.Axis{input.YAxis}

	return opt
}

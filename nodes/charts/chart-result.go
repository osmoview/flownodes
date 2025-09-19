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
	opt.Legend.Show = echarts.TRUE

	for _, t := range input.XAxis {
		opt.XAxis.Data = append(opt.XAxis.Data, t.Format(echarts.TimeLayoutDateTime))
	}

	opt.YAxis = input.YAxis

	return opt
}

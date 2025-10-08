package charts

import (
	"github.com/osmoview/echarts"
	// "github.com/osmoview/flownodes"
)

type ChartYAxis struct {
	Min any
	Max any
}

type ChartYAxisOutput struct {
	YAxis echarts.Axis
}

func (node *ChartYAxis) Run() ChartYAxisOutput {
	return ChartYAxisOutput{YAxis: echarts.Axis{
		Min:         node.Min,
		Max:         node.Max,
		BoundaryGap: [2]any{"10%", "10%"},
		SplitLine: echarts.SplitLine{
			LineStyle: echarts.LineStyle{
				Color:   "currentColor",
				Type:    "dashed",
				Opacity: 0.3,
			},
		},
	}}
}

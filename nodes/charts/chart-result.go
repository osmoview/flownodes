package charts

import "time"

type ChartResult struct{}

type ChartResultData struct {
	Series []ChartSeries
	XAxis  []time.Time
}

func (t *ChartResult) Run(data ChartResultData) ChartResultData {
	return data
}

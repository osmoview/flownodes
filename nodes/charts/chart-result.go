package charts

type ChartResult struct{}

type ChartResultData struct {
	Series []ChartSeries
	XAxis  []any
}

func (t *ChartResult) Run(data ChartResultData) ChartResultData {
	return data
}

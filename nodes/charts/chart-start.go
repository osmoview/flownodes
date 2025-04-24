package charts

import "time"

type ChartStart struct{}

type ChartStartData struct {
	Filters  ChartFilters
	Duration time.Duration
	Interval time.Duration
}

func (t *ChartStart) Run(data ChartStartData) ChartStartData {
	return data
}

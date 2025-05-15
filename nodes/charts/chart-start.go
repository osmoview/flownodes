package charts

import (
	"time"
)

type ChartStart struct{}

type ChartStartData struct {
	Duration time.Duration
	Interval time.Duration
}

func (ChartStart) Run(data ChartStartData) ChartStartData {
	return data
}

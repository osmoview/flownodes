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
	if data.Duration == 0 {
		data.Duration = 24 * time.Hour
	}
	if data.Interval == 0 {
		data.Interval = time.Hour
	}
	return data
}

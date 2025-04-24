package charts

import "time"

const xaxisTimeFull = "2006-01-02 15:04"
const xaxisTimeTime = "15:04"
const xaxisTimeDate = "2006-01-02"

type ChartXAxis struct {
}

type ChartXAxisInput struct {
	Duration time.Duration
	Interval time.Duration
}

type ChartXAxisOutput struct {
	XAxis []any
}

func (node *ChartXAxis) Run(in ChartXAxisInput) (out ChartXAxisOutput) {
	// prepare time range and interval
	timeto := time.Now()
	timeat := timeto.Add(-in.Duration)

	layout := xaxisTimeFull
	// 2 days
	if in.Duration.Seconds() < 86400 {
		layout = xaxisTimeTime
	} else if in.Duration.Seconds() > 86400*7 {
		layout = xaxisTimeDate
	}

	// fill xaxis time values
	for x := timeat; x.Before(timeto); x = x.Add(in.Interval) {
		out.XAxis = append(out.XAxis, x.Format(layout))
	}

	return
}

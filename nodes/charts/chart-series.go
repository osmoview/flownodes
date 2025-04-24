package charts

import (
	"slices"
	"sort"
	"time"

	"github.com/osmoview/flowdata"
	"github.com/osmoview/flownodes"
)

const (
	SplitByClass      = "class"
	SplitByStatus     = "status"
	SplitByDataSource = "datasource"
	SplitByFlow       = "flow"
)

type ChartBuildSeries struct {
	storage flownodes.Storage
	db      flownodes.Database

	Type  string
	Split string
}

type ChartBuildSeriesInput struct {
	Filters  ChartFilters
	Duration time.Duration
	Interval time.Duration
}

type ChartBuildSeriesOutput struct {
	Series []ChartSeries
}

func (node *ChartBuildSeries) Init(storage flownodes.Storage, db flownodes.Database) {
	node.storage = storage
	node.db = db
}

func (node *ChartBuildSeries) Run(in ChartBuildSeriesInput) (out ChartBuildSeriesOutput, err error) {
	timeto := time.Now()
	timeat := timeto.Add(-in.Duration)
	size := int(in.Duration / in.Interval)
	series := map[string]ChartSeries{}

	for i := 0; i < size; i++ {
		t := timeat.Add(time.Duration(i) * in.Interval)
		if err = node.loadEvents(t, t.Add(in.Interval), in.Filters, func(e flowdata.Event) bool {
			switch node.Split {
			case SplitByClass:
				for _, o := range e.Objects {
					if len(in.Filters.Classes) == 0 || slices.Contains(in.Filters.Classes, o.Type) {
						node.setSeries(series, o.Type, i, size, 1)
					}
				}
			case SplitByStatus:
				node.setSeries(series, e.Status, i, size, 1)
			case SplitByDataSource:
				ds, err := node.storage.GetDataSource(e.DataSource)
				if err == nil {
					node.setSeries(series, ds.GetName(), i, size, 1)
				}
			case SplitByFlow:
				flw, err := node.storage.GetFlow(e.Flow)
				if err == nil {
					node.setSeries(series, flw.GetName(), i, size, 1)
				}
			default:
				node.setSeries(series, "", i, size, float64(len(e.Objects)))
			}

			return true
		}); err != nil {
			return out, err
		}
	}

	for _, s := range series {
		out.Series = append(out.Series, s)
	}

	sort.Slice(out.Series, func(i, j int) bool {
		return out.Series[i].Name < out.Series[j].Name
	})

	return
}

func (node *ChartBuildSeries) loadEvents(timeat, timeto time.Time, filters ChartFilters, f func(e flowdata.Event) bool) error {
	// filter := events.GetEventsQuery{
	// 	TimeAt:      timeat.String(),
	// 	TimeTo:      timeto.String(),
	// 	DataSources: filters.DataSources,
	// 	Flows:       filters.Flows,
	// 	Status:      filters.Status,
	// }

	node.db.Table(flownodes.TableScheme{Table: "events"})
	// return events.RangeEvents(filter, f)
	return nil
}

func (node *ChartBuildSeries) setSeries(series map[string]ChartSeries, name string, index, size int, value float64) {
	data, ok := series[name]
	if !ok {
		data = ChartSeries{
			Name: name,
			Type: node.Type,
			Data: make([]float64, size),
		}
		series[name] = data
	}

	data.Data[index] += value
}

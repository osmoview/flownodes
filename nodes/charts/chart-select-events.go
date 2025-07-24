package charts

import (
	"log"
	"strings"
	"time"

	"github.com/osmoview/flowdata"
	"github.com/osmoview/flownodes"
)

type ChartSelectEvents struct {
	DataSources []string
	Flows       []string
	Classes     []string
	Status      string

	table flownodes.DatabaseTable
}

type ChartSelectEventsInput struct {
	Duration time.Duration
}

type ChartSelectEventsOutput struct {
	Events []*flowdata.Event
}

func (node *ChartSelectEvents) Init(db flownodes.Database) (err error) {
	node.table, err = db.Table(flownodes.TableScheme{Table: "events"})
	return err
}

func (node *ChartSelectEvents) Run(in ChartSelectEventsInput) (out ChartSelectEventsOutput, err error) {
	err = node.table.Range(node.buildFilter(in), func(row flownodes.DatabaseTableRow) bool {
		var e flowdata.Event
		if err = row.BindTo(&e); err != nil {
			log.Println(err)
			return true
		}

		// check by status
		if len(node.Status) > 0 && !strings.Contains(e.Status, node.Status) {
			return true
		}

		// check by objects
		if len(node.Classes) > 0 && !node.HasObject(e.Objects, node.Classes) {
			return true
		}

		out.Events = append(out.Events, &e)

		return true
	})

	return
}

func (node *ChartSelectEvents) buildFilter(in ChartSelectEventsInput) (filter map[string]any) {
	filter = make(map[string]any)

	if len(node.DataSources) > 0 {
		filter["datasource"] = node.DataSources
	}

	if len(node.Flows) > 0 {
		filter["flow"] = node.Flows
	}

	filter["created_at >= ?"] = time.Now().Add(-1 * in.Duration)

	return
}

func (node *ChartSelectEvents) HasObject(objects []*flowdata.Object, classes []string) bool {
	for _, obj := range objects {
		for _, c := range classes {
			if obj.Type == c {
				return true
			}
		}
	}

	return false
}

package charts

type ChartFilters struct {
	DataSources []string
	Flows       []string
	Status      string
	Classes     []string
}

type ChartSeries struct {
	Name string
	Type string
	Data []float64
}

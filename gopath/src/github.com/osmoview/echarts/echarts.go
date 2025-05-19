package echarts

const TimeLayoutDateTime = "2006-01-02 15:04"
const TimeLayoutTime = "15:04"
const TimeLayoutDate = "2006-01-02"

var (
	_true  = true
	_false = false

	// just for a pointer
	TRUE  = &_true
	FALSE = &_false
)

func NewEcharts() (eopt Echarts) {
	// prepare echarts options
	eopt.Title = Title{
		Show: FALSE,
	}

	eopt.Tooltip.Trigger = "axis"

	eopt.Grid = Grid{
		Show:         FALSE,
		Left:         "2%",
		Top:          "3%",
		Right:        "2%",
		Bottom:       "3%",
		ContainLabel: TRUE,
	}

	eopt.Legend = Legend{
		Show:  FALSE,
		Top:   "0px",
		Right: "2%",
		Icon:  "roundRect",
		TextStyle: TextStyle{
			Color:      "currentColor",
			FontSize:   14,
			FontFamily: "Mulish",
		},
	}

	eopt.TextStyle = TextStyle{
		Color:           "currentColor", //"#fff",
		FontSize:        16,
		FontFamily:      "Mulish",
		TextBorderWidth: 0,
	}

	eopt.YAxis.Min = 0

	return eopt
}

type Echarts struct {
	Title     Title     `json:"title"`
	Legend    Legend    `json:"legend"`
	Tooltip   Tooltip   `json:"tooltip"`
	Grid      Grid      `json:"grid"`
	XAxis     Axis      `json:"xAxis"`
	YAxis     Axis      `json:"yAxis"`
	Series    []Series  `json:"series"`
	TextStyle TextStyle `json:"textStyle"`
	BlendMode string    `json:"blendMode,omitempty"`
}

type Title struct {
	Text string `json:"text,omitempty"`
	Show *bool  `json:"show,omitempty"`
	Link string `json:"link,omitempty"`
}

type Tooltip struct {
	Trigger string `json:"trigger,omitempty"`
}

type Legend struct {
	Show   *bool  `json:"show,omitempty"`
	Left   string `json:"left,omitempty"`
	Top    string `json:"top,omitempty"`
	Right  string `json:"right,omitempty"`
	Bottom string `json:"bottom,omitempty"`
	//...
	Icon      string    `json:"icon,omitempty"`
	TextStyle TextStyle `json:"textStyle"`
}

type Grid struct {
	Left         string `json:"left,omitempty"`
	Top          string `json:"top,omitempty"`
	Right        string `json:"right,omitempty"`
	Bottom       string `json:"bottom,omitempty"`
	Show         *bool  `json:"show,omitempty"`
	ContainLabel *bool  `json:"containLabel,omitempty"`
}

type Axis struct {
	Show        *bool     `json:"show,omitempty"`
	Type        string    `json:"type,omitempty"`
	Data        []any     `json:"data,omitempty"`
	AxisLine    AxisLine  `json:"axisLine"`
	SplitLine   SplitLine `json:"splitLine"`
	AxisTick    AxisTick  `json:"axisTick"`
	BoundaryGap [2]any    `json:"boundaryGap,omitempty"`
	Min         any       `json:"min,omitempty"`
	Max         any       `json:"max,omitempty"`
}

type AxisLine struct {
	Show      *bool `json:"show,omitempty"`
	LineStyle LineStyle
}

type AxisTick struct {
	Show   *bool `json:"show,omitempty"`
	Inside *bool `json:"inside,omitempty"`
	Length int   `json:"length,omitempty"`
}

type SplitLine struct {
	Show      *bool     `json:"show,omitempty"`
	LineStyle LineStyle `json:"lineStyle"`
}

type Series struct {
	Type       string      `json:"type"`
	Name       string      `json:"name"`
	Data       []float64   `json:"data,omitempty"`
	LineStyle  LineStyle   `json:"lineStyle"`
	ItemStyle  ItemStyle   `json:"itemStyle"`
	Smooth     any         `json:"smooth,omitempty"`
	ShowSymbol *bool       `json:"showSymbol,omitempty"`
	Label      SeriesLabel `json:"label,omitempty"`
	Stack      string      `json:"stack,omitempty"`
}

type SeriesLabel struct {
	Show     *bool  `json:"show,omitempty"`
	Position string `json:"position,omitempty"`
	TextStyle
}

type LineStyle struct {
	Color   string  `json:"color,omitempty"`
	Width   float32 `json:"width,omitempty"`
	Type    string  `json:"type,omitempty"`
	Opacity float32 `json:"opacity,omitempty"`
}

type ItemStyle struct {
	BorderRadius any `json:"borderRadius,omitempty"`
}

type TextStyle struct {
	Color           string  `json:"color,omitempty"`
	FontSize        float32 `json:"fontSize,omitempty"`
	LineHeight      float32 `json:"lineHeight,omitempty"`
	FontFamily      string  `json:"fontFamily,omitempty"`
	TextBorderWidth any     `json:"textBorderWidth,omitempty"`
	TextBorderColor string  `json:"textBorderColor,omitempty"`
	TextShadowColor string  `json:"textShadowColor,omitempty"`
}

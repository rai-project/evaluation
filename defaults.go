package evaluation

const (
	GoldRatio = 1.618033989
)

var (
	DefaultTitleFontSize  = 25
	DefaultSeriesFontSize = 14
	DefaultLegendFontSize = 14
	DefaultBarPlotWidth   = 900
	DefaultBarPlotHeight  = int(float64(DefaultBarPlotWidth) / GoldRatio)
	DefaultBoxPlotWidth   = 900
	DefaultBoxPlotHeight  = int(float64(DefaultBoxPlotWidth) / GoldRatio)
	DefaultPiePlotWidth   = 900
	DefaultPiePlotHeight  = 500
)

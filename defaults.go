package evaluation

const (
	GoldRatio = 1.618033989
)

var (
	DefaultShowTitle          = true
	DefaultAssetHost          = `http://chenjiandongx.com/go-echarts-assets/assets/`
	DefaultTitleFontSize      = 25
	DefaultSeriesFontSize     = 14
	DefaultLegendFontSize     = 14
	DefaultBarPlotAspectRatio = 3.0
	DefaultBarPlotWidth       = 900
	DefaultBarPlotHeight      = int(float64(DefaultBarPlotWidth) / DefaultBarPlotAspectRatio)
	DefaultBoxPlotAspectRatio = 3.0
	DefaultBoxPlotWidth       = 900
	DefaultBoxPlotHeight      = int(float64(DefaultBoxPlotWidth) / DefaultBoxPlotAspectRatio)
	DefaultPiePlotWidth       = 900
	DefaultPiePlotHeight      = 500
)

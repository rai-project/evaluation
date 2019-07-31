package evaluation

const (
	GoldRatio = 1.618033989
)

var (
	DefaultShowTitle          = true
	DefaultAssetHost          = `https://s3.amazonaws.com/store.carml.org/model_analysis_2019/assets/`
	DefaultTitleFontSize      = 18
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

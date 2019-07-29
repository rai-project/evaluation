package evaluation

import (
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/rai-project/go-echarts/charts"
	"github.com/rai-project/utils/browser"
)

type PlotNamed interface {
	PlotName() string
}

type BarPlotter interface {
	PlotNamed
	BarPlot() *charts.Bar
	BarPlotAdd(*charts.Bar) *charts.Bar
	WriteBarPlot(string) error
	OpenBarPlot() error
}

type BoxPlotter interface {
	PlotNamed
	BoxPlot() *charts.BoxPlot
	BoxPlotAdd(*charts.BoxPlot) *charts.BoxPlot
	WriteBoxPlot(string) error
	OpenBoxPlot() error
}

type PiePlotter interface {
	PlotNamed
	PiePlot() *charts.Pie
	PiePlotAdd(*charts.Pie) *charts.Pie
	WritePiePlot(string) error
	OpenPiePlot() error
}

func writeBarPlot(o BarPlotter, filepath string) error {
	bar := o.BarPlot()

	if DefaultShowTitle {
		bar.SetGlobalOptions(
			charts.TitleOpts{
				Title: o.PlotName(),
				Right: "center",
				Top:   "top",
				TitleStyle: charts.TextStyleOpts{
					FontSize: DefaultTitleFontSize,
				},
			})
	}

	bar.SetGlobalOptions(
		charts.LegendOpts{
			Right: "right",
			Top:   "middle",
			TextStyle: charts.TextStyleOpts{
				FontSize: DefaultLegendFontSize,
			},
		},
		charts.ToolboxOpts{Show: true, TBFeature: charts.TBFeature{SaveAsImage: charts.SaveAsImage{PixelRatio: 5}}},
		charts.InitOpts{
			AssetsHost: DefaultAssetHost,
			Theme:      charts.ThemeType.Shine,
			Width:      fmt.Sprintf("%vpx", DefaultBarPlotWidth),
			Height:     fmt.Sprintf("%vpx", DefaultBarPlotHeight),
		},
		// charts.DataZoomOpts{XAxisIndex: []int{0}, Start: 50, End: 100},
	)
	os.MkdirAll(path.Dir(filepath), os.ModePerm)
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	err = bar.Render(f)
	if err != nil {
		return err
	}
	return nil
}

func writeBoxPlot(o BoxPlotter, filepath string) error {
	box := o.BoxPlot()

	if DefaultShowTitle {
		box.SetGlobalOptions(
			charts.TitleOpts{
				Title: o.PlotName(),
				Right: "center",
				Top:   "top",
				TitleStyle: charts.TextStyleOpts{
					FontSize: DefaultTitleFontSize,
				},
			})
	}

	box.SetGlobalOptions(
		charts.LegendOpts{
			Right: "80%",
			TextStyle: charts.TextStyleOpts{
				FontSize: DefaultLegendFontSize,
			},
		},
		charts.ToolboxOpts{Show: true, TBFeature: charts.TBFeature{SaveAsImage: charts.SaveAsImage{PixelRatio: 5}}},
		charts.InitOpts{
			AssetsHost: DefaultAssetHost,
			Theme:      charts.ThemeType.Shine,
			Width:      fmt.Sprintf("%vpx", DefaultBarPlotWidth),
			Height:     fmt.Sprintf("%vpx", DefaultBarPlotHeight),
		},
	)
	os.MkdirAll(path.Dir(filepath), os.ModePerm)
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	err = box.Render(f)
	if err != nil {
		return err
	}
	return nil
}

func writePiePlot(o PiePlotter, filepath string) error {
	pie := o.PiePlot()

	if DefaultShowTitle {
		pie.SetGlobalOptions(
			charts.TitleOpts{
				Title: o.PlotName(),
				Right: "center",
				Top:   "bottom",
				TitleStyle: charts.TextStyleOpts{
					FontSize: DefaultLegendFontSize,
				},
			})
	}

	pie.SetGlobalOptions(
		charts.LegendOpts{
			Right: "right",
			Top:   "middle",
			TextStyle: charts.TextStyleOpts{
				FontSize: DefaultLegendFontSize,
			},
		},
		charts.ToolboxOpts{Show: true, TBFeature: charts.TBFeature{SaveAsImage: charts.SaveAsImage{PixelRatio: 5}}},
		charts.InitOpts{
			AssetsHost: DefaultAssetHost,
			Width:      fmt.Sprintf("%vpx", DefaultPiePlotWidth),
			Height:     fmt.Sprintf("%vpx", DefaultPiePlotHeight),
		},
	)
	os.MkdirAll(path.Dir(filepath), os.ModePerm)
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	err = pie.Render(f)
	if err != nil {
		return err
	}
	return nil
}

func openBarPlot(o BarPlotter) error {
	filepath := TempFile("", "batchPlot_*.html")
	if filepath == "" {
		return errors.New("failed to create temporary file")
	}
	err := o.WriteBarPlot(filepath)
	if err != nil {
		return err
	}
	// defer os.Remove(filepath)
	if ok := browser.Open(filepath); !ok {
		return errors.New("failed to open browser filepath")
	}

	return nil
}

func openBoxPlot(o BoxPlotter) error {
	filepath := TempFile("", "batchPlot_*.html")
	if filepath == "" {
		return errors.New("failed to create temporary file")
	}
	err := o.WriteBoxPlot(filepath)
	if err != nil {
		return err
	}
	// defer os.Remove(filepath)
	if ok := browser.Open(filepath); !ok {
		return errors.New("failed to open browser filepath")
	}

	return nil
}

func openPiePlot(o PiePlotter) error {
	filepath := TempFile("", "batchPlot_*.html")
	if filepath == "" {
		return errors.New("failed to create temporary file")
	}
	err := o.WritePiePlot(filepath)
	if err != nil {
		return err
	}
	// defer os.Remove(filepath)
	if ok := browser.Open(filepath); !ok {
		return errors.New("failed to open browser filepath")
	}

	return nil
}

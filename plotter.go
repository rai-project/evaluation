package evaluation

import (
	"fmt"
	"os"

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

func writeBarPlot(o BarPlotter, path string) error {
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
			Right: "80%",
			TextStyle: charts.TextStyleOpts{
				FontSize: DefaultLegendFontSize,
			},
		},
		charts.ToolboxOpts{Show: true, TBFeature: charts.TBFeature{SaveAsImage: charts.SaveAsImage{PixelRatio: 5}}},
		charts.InitOpts{
			Theme:  charts.ThemeType.Shine,
			Width:  fmt.Sprintf("%vpx", DefaultBarPlotWidth),
			Height: fmt.Sprintf("%vpx", DefaultBarPlotHeight),
		},
		// charts.DataZoomOpts{XAxisIndex: []int{0}, Start: 50, End: 100},
	)
	f, err := os.Create(path)
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

func writeBoxPlot(o BoxPlotter, path string) error {
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
			Theme:  charts.ThemeType.Shine,
			Width:  fmt.Sprintf("%vpx", DefaultBarPlotWidth),
			Height: fmt.Sprintf("%vpx", DefaultBarPlotHeight),
		},
	)
	f, err := os.Create(path)
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

func writePiePlot(o PiePlotter, path string) error {
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
	)
	f, err := os.Create(path)
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
	path := TempFile("", "batchPlot_*.html")
	if path == "" {
		return errors.New("failed to create temporary file")
	}
	err := o.WriteBarPlot(path)
	if err != nil {
		return err
	}
	// defer os.Remove(path)
	if ok := browser.Open(path); !ok {
		return errors.New("failed to open browser path")
	}

	return nil
}

func openBoxPlot(o BoxPlotter) error {
	path := TempFile("", "batchPlot_*.html")
	if path == "" {
		return errors.New("failed to create temporary file")
	}
	err := o.WriteBoxPlot(path)
	if err != nil {
		return err
	}
	// defer os.Remove(path)
	if ok := browser.Open(path); !ok {
		return errors.New("failed to open browser path")
	}

	return nil
}

func openPiePlot(o PiePlotter) error {
	path := TempFile("", "batchPlot_*.html")
	if path == "" {
		return errors.New("failed to create temporary file")
	}
	err := o.WritePiePlot(path)
	if err != nil {
		return err
	}
	// defer os.Remove(path)
	if ok := browser.Open(path); !ok {
		return errors.New("failed to open browser path")
	}

	return nil
}

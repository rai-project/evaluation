package evaluation

import (
	"os"

	"github.com/pkg/errors"
	"github.com/rai-project/go-echarts/charts"
	"github.com/rai-project/utils/browser"
)

type Plotter interface {
	Name() string
	BarPlot(string) *charts.Bar
	BarPlotAdd(*charts.Bar) *charts.Bar
	WritePlot(string) error
	OpenPlot() error
}

func writePlot(o Plotter, path string) error {
	bar := o.BarPlot(o.Name())
	bar.SetGlobalOptions(
		charts.TitleOpts{Right: "40%"},
		charts.LegendOpts{Right: "80%"},
		charts.ToolboxOpts{Show: true},
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

func openPlot(o Plotter) error {
	path := TempFile("", "batchPlot_*.html")
	if path == "" {
		return errors.New("failed to create temporary file")
	}
	err := o.WritePlot(path)
	if err != nil {
		return err
	}
	// defer os.Remove(path)
	if ok := browser.Open(path); !ok {
		return errors.New("failed to open browser path")
	}

	return nil
}

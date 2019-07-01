package evaluation

import (
	"os"

	"github.com/pkg/errors"
	"github.com/rai-project/go-echarts/charts"
	"github.com/rai-project/utils/browser"
)

type Named interface {
	Name() string
}

type BarPlotter interface {
	Named
	BarPlot(string) *charts.Bar
	BarPlotAdd(*charts.Bar) *charts.Bar
	WriteBarPlot(string) error
	OpenBarPlot() error
}

type BoxPlotter interface {
	Named
	BoxPlot(string) *charts.BoxPlot
	BoxPlotAdd(*charts.BoxPlot) *charts.BoxPlot
	WriteBoxPlot(string) error
	OpenBoxPlot() error
}

func writeBarPlot(o BarPlotter, path string) error {
	bar := o.BarPlot(o.Name())
	bar.SetGlobalOptions(
		charts.TitleOpts{Right: "40%"},
		charts.LegendOpts{Right: "80%"},
		charts.ToolboxOpts{Show: true},
		charts.InitOpts{Theme: charts.ThemeType.Shine},
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
	box := o.BoxPlot(o.Name())
	box.SetGlobalOptions(
		charts.TitleOpts{Right: "40%"},
		charts.LegendOpts{Right: "80%"},
		charts.ToolboxOpts{Show: true},
		charts.InitOpts{Theme: charts.ThemeType.Shine},
		// charts.DataZoomOpts{XAxisIndex: []int{0}, Start: 50, End: 100},
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

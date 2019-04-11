package plotting

import (
	"net/http"
	"os"

	"github.com/chenjiandongx/go-echarts/charts"
	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"github.com/rai-project/evaluation/plotting/browser"
	"github.com/spf13/cast"
)

type batchPlot struct {
	Name      string
	Durations []batchDurationSummary
	Options   *Options
}

type batchDurationSummary struct {
	BatchSize int
	Duration  int64
}

func NewBatchPlot(name string, os ...OptionModifier) (*batchPlot, error) {
	os = append(os, Option.BatchSize(0))
	opts := NewOptions(os...)

	traceFilePaths, err := modelTraceFiles(opts)
	if err != nil {
		return nil, err
	}

	durations := []batchDurationSummary{}
	for _, path := range traceFilePaths {
		tr, err := ReadTraceFile(path)
		if err != nil {
			if opts.ignoreReadErrors {
				continue
			}
			return nil, err
		}
		cPredictSpan, err := findCPredict(tr)
		if err != nil {
			return nil, err
		}
		batchSize, err := findBatchSize(tr)
		if err != nil {
			return nil, err
		}
		durations = append(durations,
			batchDurationSummary{
				BatchSize: batchSize,
				Duration:  int64(cPredictSpan.Duration),
			},
		)
	}
	return &batchPlot{
		Name:      name,
		Durations: durations,
		Options:   opts,
	}, nil
}

func (o batchPlot) BarPlot() *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.TitleOpts{Title: "xxx"},
		charts.ToolboxOpts{Show: true},
	)
	labels := make([]string, len(o.Durations))
	for ii, elem := range o.Durations {
		labels[ii] = cast.ToString(elem.BatchSize)
	}
	durations := make([]int64, len(o.Durations))
	for ii, elem := range o.Durations {
		durations[ii] = elem.Duration
	}
	bar.AddXAxis(labels).AddYAxis(o.Name, durations)
	return bar
}

func (o batchPlot) Write(path string) error {
	bar := o.BarPlot()
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

func (o batchPlot) Open() error {
	path := tempFile("", "batchPlot_*.html")
	if path == "" {
		return errors.New("failed to create temporary file")
	}
	err := o.Write(path)
	if err != nil {
		return err
	}
	// defer os.Remove(path)
	if ok := browser.Open(path); !ok {
		return errors.New("failed to open browser path")
	}

	return nil
}

func (o batchPlot) Handler(w http.ResponseWriter, _ *http.Request) {
	bar := o.BarPlot()
	bar.Render(w)
}

func (o batchPlot) Header() []string {
	return batchDurationSummary{}.Header()
}

func (o batchPlot) Rows() [][]string {
	res := make([][]string, len(o.Durations))
	for ii, dur := range o.Durations {
		res[ii] = dur.Row()
	}
	return res
}

func (o batchDurationSummary) Header() []string {
	res := []string{}
	s := structs.New(&o)
	for _, field := range s.Fields() {
		res = append(res, field.Name())
	}
	return res
}

func (o batchDurationSummary) Row() []string {
	res := []string{}
	s := structs.New(&o)
	for _, field := range s.Fields() {
		res = append(res, cast.ToString(field.Value()))
	}
	return res
}

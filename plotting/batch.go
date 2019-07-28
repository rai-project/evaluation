package plotting

import (
	"net/http"
	"os"
	"sort"

	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"github.com/rai-project/evaluation/writer"
	"github.com/rai-project/go-echarts/charts"
	"github.com/rai-project/utils/browser"
	"github.com/spf13/cast"
)

//go:generate collections -file $GOFILE

type batchPlot struct {
	Name      string
	Durations []batchDurationSummary
	Options   *Options
}

type batchDurationSummary struct {
	BatchSize int
	ModelName string
	Duration  int64
}

type batchDurationSummaries []batchDurationSummary

func (p batchDurationSummaries) Len() int           { return len(p) }
func (p batchDurationSummaries) Less(i, j int) bool { return p[i].BatchSize < p[j].BatchSize }
func (p batchDurationSummaries) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

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
		modelName, err := findModelName(tr)
		if err != nil {
			return nil, err
		}
		durations = append(durations,
			batchDurationSummary{
				BatchSize: batchSize,
				ModelName: modelName,
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

func (o batchPlot) BarPlotAdd(bar *charts.Bar) *charts.Bar {
	labels := []int{}
	for _, elem := range o.Durations {
		// batchSize := int(math.Log2(float64(elem.BatchSize)))
		batchSize := elem.BatchSize
		if contains(labels, batchSize) {
			continue
		}
		labels = append(labels, batchSize)
	}

	sort.Sort(sort.IntSlice(labels))

	// modelDurations := orderedmap.New()
	// for _, elem := range o.Durations {
	// 	var val []int64
	// 	if e, ok := modelDurations.Get(elem.ModelName); ok {
	// 		val = e.([]int64)
	// 	} else {
	// 		val = []int64{}
	// 	}
	// 	val = append(val, elem.Duration)
	// 	modelDurations.Set(elem.ModelName, val)
	// }

	modelDurations := map[string]batchDurationSummaries{}
	for _, elem := range o.Durations {
		var val []batchDurationSummary
		if e, ok := modelDurations[elem.ModelName]; ok {
			val = e
		} else {
			val = []batchDurationSummary{}
		}
		val = append(val, elem)
		modelDurations[elem.ModelName] = batchDurationSummaries(val)
	}
	strLabels := make([]string, len(labels))
	for ii, label := range labels {
		strLabels[ii] = cast.ToString(label)
	}
	bar.AddXAxis(strLabels)
	for key, durations := range modelDurations {
		vals := make([]int64, len(labels))
		for ii, label := range labels {
			for _, duration := range durations {
				if duration.BatchSize == label {
					vals[ii] = duration.Duration / 1000
					break
				}
			}
		}
		// pp.Println(key, " --", vals)
		bar.AddYAxis(key, vals)
	}
	bar.SetSeriesOptions(
		charts.LabelTextOpts{Show: false},
		charts.TextStyleOpts{FontSize: DefaultSeriesFontSize},
	)
	bar.SetGlobalOptions(
		charts.XAxisOpts{Name: "BatchSize", Show: false, AxisLabel: charts.LabelTextOpts{Show: true}},
		charts.YAxisOpts{Name: "Latency(ms)"},
	)
	return bar
}

func (o batchPlot) BarPlot() *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.ToolboxOpts{Show: true, TBFeature{SaveAsImage: struct{pixelRatio: 5}}},
	)
	bar = o.BarPlotAdd(bar)
	return bar
}

func (o batchPlot) Write(path string) error {
	bar := o.BarPlot(o.Name)
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
	bar := o.BarPlot(o.Name)
	bar.Render(w)
}

func (o batchPlot) Header(opts ...writer.Option) []string {
	return batchDurationSummary{}.Header(opts...)
}

func (o batchPlot) Rows(opts ...writer.Option) [][]string {
	res := make([][]string, len(o.Durations))
	for ii, dur := range o.Durations {
		res[ii] = dur.Row(opts...)
	}
	return res
}

func (o batchDurationSummary) Header(opts ...writer.Option) []string {
	res := []string{}
	s := structs.New(&o)
	for _, field := range s.Fields() {
		res = append(res, field.Name())
	}
	return res
}

func (o batchDurationSummary) Row(opts ...writer.Option) []string {
	res := []string{}
	s := structs.New(&o)
	for _, field := range s.Fields() {
		res = append(res, cast.ToString(field.Value()))
	}
	return res
}

package evaluation

import (
	"errors"
	"fmt"

	"github.com/rai-project/evaluation/writer"
	"github.com/rai-project/go-echarts/charts"
	"github.com/rai-project/tracer"
	"github.com/spf13/cast"
)

//easyjson:json
type SummaryModelInformation struct {
	SummaryBase `json:",inline,omitempty"`
	Durations   []int64 `json:"durations,omitempty"`
	Duration    float64 `json:"duration,omitempty"`
	Latency     float64 `json:"latency,omitempty"`
	Throughput  float64 `json:"throughput,omitempty"`
}

type SummaryModelInformations []SummaryModelInformation

func (p SummaryModelInformations) Len() int { return len(p) }
func (p SummaryModelInformations) Less(i, j int) bool {
	return p[i].BatchSize < p[j].BatchSize
}
func (p SummaryModelInformations) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type SummaryModelLatencyInformations SummaryModelInformations

type SummaryModelThroughputInformations SummaryModelInformations

func (SummaryModelInformation) Header(opts ...writer.Option) []string {
	extra := []string{
		"duration (us)",
		"latency (ms)",
		"throughput (input/s)",
		// "durations (us)",
	}
	return append(SummaryBase{}.Header(opts...), extra...)
}

func (s SummaryModelInformation) Row(opts ...writer.Option) []string {
	extra := []string{
		fmt.Sprintf("%.2f", s.Duration),
		fmt.Sprintf("%.2f", s.Latency),
		fmt.Sprintf("%.2f", s.Throughput),
		// strings.Join(int64SliceToStringSlice(s.Durations), ","),
	}
	return append(s.SummaryBase.Row(opts...), extra...)
}

func (es Evaluations) SummaryModelInformations(perfCol *PerformanceCollection) (SummaryModelInformations, error) {
	summary := SummaryModelInformations{}
	if len(es) == 0 {
		return summary, errors.New("no evaluation is found in the database")
	}
	groupedEvals := es.GroupByBatchSize()

	for _, evals := range groupedEvals {
		spans, err := evals.GetSpansFromPerformanceCollection(perfCol)
		if err != nil {
			return summary, err
		}
		if len(spans) == 0 {
			return summary, errors.New("no span is found for the evaluation")
		}

		cPredictSpans := spans.FilterByOperationNameAndEvalTraceLevel("c_predict", tracer.MODEL_TRACE.String())

		durations := []int64{}
		for _, span := range cPredictSpans {
			durations = append(durations, cast.ToInt64(span.Duration))
		}
		duration := TrimmedMeanInt64Slice(durations, DefaultTrimmedMeanFraction)
		base := evals[0].summaryBase()
		batchSize := base.BatchSize
		if duration == 0 {
			continue
		}
		latency := duration / float64(batchSize*1000)
		summary = append(summary, SummaryModelInformation{
			SummaryBase: base,
			Durations:   durations,
			Duration:    duration,
			Throughput:  float64(1000) / latency,
			Latency:     latency,
		})
	}
	return summary, nil
}

func (o SummaryModelThroughputInformations) PlotName() string {
	if len(o) == 0 {
		return ""
	}
	return o[0].ModelName + `
  Throughput`
}

func (o SummaryModelLatencyInformations) PlotName() string {
	if len(o) == 0 {
		return ""
	}
	return o[0].ModelName + `
  Batch Latency`
}

func (o SummaryModelThroughputInformations) BarPlot() *charts.Bar {
	bar := charts.NewBar()
	bar = o.BarPlotAdd(bar)
	return bar
}

func (o SummaryModelLatencyInformations) BarPlot() *charts.Bar {
	bar := charts.NewBar()
	bar = o.BarPlotAdd(bar)
	return bar
}

type SummaryModelInformationsSelector func(elem SummaryModelInformation) float64

func (o SummaryModelInformations) barPlotAdd(bar *charts.Bar, elemSelector SummaryModelInformationsSelector) *charts.Bar {
	labels := []string{}
	for _, elem := range o {
		labels = append(labels, cast.ToString(elem.BatchSize))
	}
	bar.AddXAxis(labels)

	data := make([]float64, len(o))
	for ii, elem := range o {
		data[ii] = elemSelector(elem)
	}
	bar.AddYAxis("", data)

	bar.SetSeriesOptions(
		charts.LabelTextOpts{Show: false},
		charts.TextStyleOpts{FontSize: DefaultSeriesFontSize},
	)

	bar.SetGlobalOptions(
		charts.XAxisOpts{Name: "Batch Size", Show: false, AxisLabel: charts.LabelTextOpts{Show: true}},
	)
	return bar
}

func (o SummaryModelThroughputInformations) BarPlotAdd(bar0 *charts.Bar) *charts.Bar {
	bar := SummaryModelInformations(o).barPlotAdd(bar0, func(elem SummaryModelInformation) float64 {
		return elem.Throughput
	})
	bar.SetGlobalOptions(
		charts.YAxisOpts{Name: "Throughput (inputs/second)"},
	)
	return bar
}

func (o SummaryModelLatencyInformations) BarPlotAdd(bar0 *charts.Bar) *charts.Bar {
	bar := SummaryModelInformations(o).barPlotAdd(bar0, func(elem SummaryModelInformation) float64 {
		return float64(elem.Duration) / float64(1000)
	})
	bar.SetGlobalOptions(
		charts.YAxisOpts{Name: "Batch Latency (ms)"},
	)
	return bar
}

func (o SummaryModelThroughputInformations) WriteBarPlot(path string) error {
	return writeBarPlot(o, path)
}

func (o SummaryModelLatencyInformations) WriteBarPlot(path string) error {
	return writeBarPlot(o, path)
}

func (o SummaryModelThroughputInformations) OpenBarPlot() error {
	return openBarPlot(o)
}

func (o SummaryModelLatencyInformations) OpenBarPlot() error {
	return openBarPlot(o)
}

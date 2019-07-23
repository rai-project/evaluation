package evaluation

import (
	"errors"
	"time"

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
		cast.ToString(s.Duration),
		cast.ToString(s.Latency),
		cast.ToString(s.Throughput),
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
		duration := time.Duration(TrimmedMeanInt64Slice(durations, DefaultTrimmedMeanFraction))
		base := evals[0].summaryBase()
		batchSize := base.BatchSize
		if duration == 0 {
			continue
		}
		latency := float64(duration) / float64(time.Duration(batchSize)*time.Microsecond)
		summary = append(summary, SummaryModelInformation{
			SummaryBase: base,
			Durations:   durations,
			Duration:    float64(duration),
			Throughput:  1 / latency,
			Latency:     latency,
		})
	}
	return summary, nil
}

func (o SummaryModelInformations) PlotName() string {
	if len(o) == 0 {
		return ""
	}
	return o[0].ModelName + " Throughput"
}

func (o SummaryModelInformations) BarPlot(title string) *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.TitleOpts{Title: title},
		charts.ToolboxOpts{Show: true},
	)
	bar = o.BarPlotAdd(bar)
	return bar
}

func (o SummaryModelInformations) BarPlotAdd(bar *charts.Bar) *charts.Bar {
	labels := []string{}
	for _, elem := range o {
		labels = append(labels, cast.ToString(elem.BatchSize))
	}
	bar.AddXAxis(labels)

	data := make([]float64, len(o))
	for ii, elem := range o {
		data[ii] = elem.Throughput
	}
	bar.AddYAxis("", data)

	bar.SetSeriesOptions(charts.LabelTextOpts{Show: false})
	bar.SetGlobalOptions(
		charts.XAxisOpts{Name: "Batch Size", Show: false, AxisLabel: charts.LabelTextOpts{Show: true}},
		charts.YAxisOpts{Name: "Throughput (inputs/second)"},
	)
	return bar
}

func (o SummaryModelInformations) WriteBarPlot(path string) error {
	return writeBarPlot(o, path)
}

func (o SummaryModelInformations) OpenBarPlot() error {
	return openBarPlot(o)
}

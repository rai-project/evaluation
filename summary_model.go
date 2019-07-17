package evaluation

import (
	"errors"
	"strings"

	"github.com/rai-project/evaluation/writer"
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

func (SummaryModelInformation) Header(opts ...writer.Option) []string {
	extra := []string{
		"durations (us)",
		"duration (us)",
		"latency (ms)",
		"throughput (input/s)",
	}
	return append(SummaryBase{}.Header(opts...), extra...)
}

func (s SummaryModelInformation) Row(opts ...writer.Option) []string {
	extra := []string{
		strings.Join(int64SliceToStringSlice(s.Durations), ","),
		cast.ToString(s.Duration),
		cast.ToString(s.Latency),
		cast.ToString(s.Throughput),
	}
	return append(s.SummaryBase.Row(opts...), extra...)
}

func summaryModelInformation(es Evaluations, spans Spans) (SummaryModelInformation, error) {
	summary := SummaryModelInformation{}
	if len(es) == 0 {
		return summary, errors.New("no evaluation is found in the database")
	}

	cPredictSpans := spans.FilterByOperationNameAndEvalTraceLevel("c_predict", tracer.MODEL_TRACE.String())
	durations := []int64{}
	for _, span := range cPredictSpans {
		durations = append(durations, cast.ToInt64(span.Duration))
	}
	duration := TrimmedMeanInt64Slice(durations, DefaultTrimmedMeanFraction)
	base := es[0].summaryBase()
	batchSize := base.BatchSize
	summary = SummaryModelInformation{
		SummaryBase: base,
		Durations:   durations,
		Duration:    duration,
		Throughput:  float64(1000000*batchSize) / duration,
		Latency:     duration / float64(batchSize*1000),
	}
	return summary, nil
}

func (es Evaluations) SummaryModelInformation(perfCol *PerformanceCollection) (SummaryModelInformation, error) {
	summary := SummaryModelInformation{}
	spans, err := es.GetSpansFromPerformanceCollection(perfCol)
	if err != nil {
		return summary, err
	}
	if len(spans) == 0 {
		return summary, errors.New("no span is found for the evaluation")
	}
	return summaryModelInformation(es, spans)
}

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

type SummaryModelInformations []SummaryModelInformation

func (SummaryModelInformation) Header(opts ...writer.Option) []string {
	extra := []string{
		"duration (us)",
		"durations (us)",
		"latency (ms)",
		"throughput (input/s)",
	}
	return append(SummaryBase{}.Header(opts...), extra...)
}

func (s SummaryModelInformation) Row(opts ...writer.Option) []string {
	extra := []string{
		cast.ToString(s.Duration),
		strings.Join(int64SliceToStringSlice(s.Durations), ","),
		cast.ToString(s.Latency),
		cast.ToString(s.Throughput),
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
		summary = append(summary, SummaryModelInformation{
			SummaryBase: base,
			Durations:   durations,
			Duration:    duration,
			Throughput:  float64(1000000*batchSize) / duration,
			Latency:     duration / float64(batchSize*1000),
		})
	}
	return summary, nil
}

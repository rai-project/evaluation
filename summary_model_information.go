package evaluation

import (
	"errors"
	"strings"

	"github.com/rai-project/evaluation/writer"
	"github.com/rai-project/tracer"
	"github.com/spf13/cast"
	model "github.com/uber/jaeger/model/json"
	db "upper.io/db.v3"
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

	cPredictSpans := spans.FilterByOperationNameAndEvalTraceLevel("c_predict", tracer.FRAMEWORK_TRACE.String())
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
	spans := []model.Span{}
	for _, e := range es {
		foundPerfs, err := perfCol.Find(db.Cond{"_id": e.PerformanceID})
		if err != nil {
			return summary, err
		}
		if len(foundPerfs) != 1 {
			return summary, errors.New("no performance is found for the evaluation")
		}
		perf := foundPerfs[0]
		spans = append(spans, perf.Spans()...)
	}
	return summaryModelInformation(es, spans)
}

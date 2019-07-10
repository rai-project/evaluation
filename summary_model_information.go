package evaluation

import (
	"errors"
	"strings"

	"github.com/rai-project/evaluation/writer"
	"github.com/rai-project/tracer"
	db "upper.io/db.v3"
)

//easyjson:json
type SummaryModelInformation struct {
	SummaryBase `json:",inline,omitempty"`
	Durations   []uint64 `json:"durations,omitempty"`
}

type SummaryModelInformations []SummaryModelInformation

func (SummaryModelInformation) Header(opts ...writer.Option) []string {
	extra := []string{
		"durations (us)",
	}
	return append(SummaryBase{}.Header(opts...), extra...)
}

func (s SummaryModelInformation) Row(opts ...writer.Option) []string {
	extra := []string{
		strings.Join(uint64SliceToStringSlice(s.Durations), ";"),
	}
	return append(s.SummaryBase.Row(opts...), extra...)
}

func (SummaryModelInformations) Header(opts ...writer.Option) []string {
	return SummaryModelInformation{}.Header(opts...)
}

func (s SummaryModelInformations) Rows(opts ...writer.Option) [][]string {
	rows := [][]string{}
	for _, e := range s {
		rows = append(rows, e.Row(opts...))
	}
	return rows
}

func (p Performance) PredictDurationInformationSummary(e Evaluation) (*SummaryModelInformation, error) {
	cPredictSpans := p.Spans().FilterByOperationNameAndEvalTraceLevel("c_predict", tracer.MODEL_TRACE.String())
	return &SummaryModelInformation{
		SummaryBase: e.summaryBase(),
		Durations:   cPredictSpans.Duration(),
	}, nil
}

func (ps Performances) PredictDurationInformationSummary(e Evaluation) ([]*SummaryModelInformation, error) {
	res := []*SummaryModelInformation{}
	for _, p := range ps {
		s, err := p.PredictDurationInformationSummary(e)
		if err != nil {
			log.WithError(err).Error("failed to get duration information summary")
			continue
		}
		res = append(res, s)
	}
	return res, nil
}

func (e Evaluation) PredictDurationInformationSummary(perfCol *PerformanceCollection) (*SummaryModelInformation, error) {
	perfs, err := perfCol.Find(db.Cond{"_id": e.PerformanceID})
	if err != nil {
		return nil, err
	}
	if len(perfs) != 1 {
		return nil, errors.New("expecting on performance output")
	}
	perf := perfs[0]
	return perf.PredictDurationInformationSummary(e)
}

func (es Evaluations) PredictDurationInformationSummary(perfCol *PerformanceCollection) (SummaryModelInformations, error) {
	res := []SummaryModelInformation{}
	for _, e := range es {
		s, err := e.PredictDurationInformationSummary(perfCol)
		if err != nil {
			log.WithError(err).Error("failed to get duration information summary")
			continue
		}
		if s == nil {
			continue
		}
		res = append(res, *s)
	}
	return res, nil
}

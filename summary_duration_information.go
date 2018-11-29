package evaluation

import (
	"errors"
	"strings"

	db "upper.io/db.v3"
)

//easyjson:json
type SummaryPredictDurationInformation struct {
	SummaryBase `json:",inline,omitempty"`
	Durations   []uint64 `json:"durations,omitempty"` // in nano seconds
}

type SummaryPredictDurationInformations []SummaryPredictDurationInformation

func (SummaryPredictDurationInformation) Header() []string {
	extra := []string{
		"durations",
	}
	return append(SummaryBase{}.Header(), extra...)
}

func (s SummaryPredictDurationInformation) Row() []string {
	extra := []string{
		strings.Join(uint64SliceToStringSlice(s.Durations), ";"),
	}
	return append(s.SummaryBase.Row(), extra...)
}

func (SummaryPredictDurationInformations) Header() []string {
	return SummaryPredictDurationInformation{}.Header()
}

func (s SummaryPredictDurationInformations) Rows() [][]string {
	rows := [][]string{}
	for _, e := range s {
		rows = append(rows, e.Row())
	}
	return rows
}

func (p Performance) PredictDurationInformationSummary(e Evaluation) (*SummaryPredictDurationInformation, error) {
	spans := p.Spans().FilterByOperationName("c_predict")

	return &SummaryPredictDurationInformation{
		SummaryBase: e.summaryBase(),
		Durations:   spans.Duration(),
	}, nil
}

func (ps Performances) PredictDurationInformationSummary(e Evaluation) ([]*SummaryPredictDurationInformation, error) {
	res := []*SummaryPredictDurationInformation{}
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

func (e Evaluation) PredictDurationInformationSummary(perfCol *PerformanceCollection) (*SummaryPredictDurationInformation, error) {
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

func (es Evaluations) PredictDurationInformationSummary(perfCol *PerformanceCollection) (SummaryPredictDurationInformations, error) {
	res := []SummaryPredictDurationInformation{}
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

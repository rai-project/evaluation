package evaluation

import (
	"errors"
	"strings"

	"github.com/rai-project/evaluation/writer"
	"github.com/spf13/cast"
	db "upper.io/db.v3"
)

//easyjson:json
type SummaryModelAccuracyInformation struct {
	SummaryBase  `json:",inline"`
	Top1Accuracy float64 `json:"top1_accuracy,omitempty"`
	Top5Accuracy float64 `json:"top5_accuracy,omitempty"`
}

type SummaryModelAccuracyInformations []SummaryModelAccuracyInformation

func (SummaryModelAccuracyInformation) Header(opts ...writer.Option) []string {
	extra := []string{
		"top1_accuracy",
		"top5_accuracy",
	}
	return append(SummaryBase{}.Header(opts...), extra...)
}

func (s SummaryModelAccuracyInformation) Row(opts ...writer.Option) []string {
	extra := []string{
		cast.ToString(s.Top1Accuracy),
		cast.ToString(s.Top5Accuracy),
	}
	return append(s.SummaryBase.Row(opts...), extra...)
}

func (s SummaryModelAccuracyInformation) key() string {
	return strings.Join(
		[]string{
			s.ModelName,
			s.ModelVersion,
			s.FrameworkName,
			s.FrameworkVersion,
			s.HostName,
			s.MachineArchitecture,
			cast.ToString(s.UsingGPU),
		},
		",",
	)
}

func (SummaryModelAccuracyInformations) Header(opts ...writer.Option) []string {
	return SummaryModelAccuracyInformation{}.Header(opts...)
}

func (s SummaryModelAccuracyInformations) Rows(opts ...writer.Option) [][]string {
	rows := [][]string{}
	for _, e := range s {
		rows = append(rows, e.Row(opts...))
	}
	return rows
}

func (s SummaryModelAccuracyInformations) Group() (SummaryModelAccuracyInformations, error) {
	groups := map[string]SummaryModelAccuracyInformations{}

	for _, v := range s {
		k := v.key()
		if _, ok := groups[k]; !ok {
			groups[k] = SummaryModelAccuracyInformations{}
		}
		groups[k] = append(groups[k], v)
	}

	res := []SummaryModelAccuracyInformation{}
	for _, v := range groups {
		if len(v) == 0 {
			log.Error("expecting more more than one input in SummaryModelAccuracyInformations")
			continue
		}
		for _, vv := range v {
			if vv.Top1Accuracy != 0 && vv.Top5Accuracy != 0 {
				res = append(res, vv)
				break
			}
		}
	}

	return res, nil
}

func (a ModelAccuracy) PredictAccuracyInformationSummary(e Evaluation) (*SummaryModelAccuracyInformation, error) {
	return &SummaryModelAccuracyInformation{
		SummaryBase:  e.summaryBase(),
		Top1Accuracy: a.Top1,
		Top5Accuracy: a.Top5,
	}, nil
}

func (as ModelAccuracies) PredictAccuracyInformationSummary(e Evaluation) ([]*SummaryModelAccuracyInformation, error) {
	res := []*SummaryModelAccuracyInformation{}
	for _, a := range as {
		s, err := a.PredictAccuracyInformationSummary(e)
		if err != nil {
			log.WithError(err).Error("failed to get accuracy information summary")
			continue
		}
		res = append(res, s)
	}
	return res, nil
}

func (e Evaluation) PredictAccuracyInformationSummary(accCol *ModelAccuracyCollection) (*SummaryModelAccuracyInformation, error) {
	accs, err := accCol.Find(db.Cond{"_id": e.ModelAccuracyID})
	if err != nil {
		return nil, err
	}
	if len(accs) != 1 {
		return nil, errors.New("expecting on model accuracy output")
	}
	acc := accs[0]
	return acc.PredictAccuracyInformationSummary(e)
}

func (es Evaluations) PredictAccuracyInformationSummary(accCol *ModelAccuracyCollection) (SummaryModelAccuracyInformations, error) {
	res := []SummaryModelAccuracyInformation{}
	for _, e := range es {
		s, err := e.PredictAccuracyInformationSummary(accCol)
		if err != nil {
			log.WithError(err).Error("failed to get accuracy information summary")
			continue
		}
		if s == nil {
			continue
		}
		res = append(res, *s)
	}
	return res, nil
}

package evaluation

import (
	"errors"
	"strings"

	"github.com/spf13/cast"
	db "upper.io/db.v3"
)

//easyjson:json
type SummaryPredictAccuracyInformation struct {
	SummaryBase  `json:",inline"`
	Top1Accuracy float64 `json:"top1_accuracy,omitempty"`
	Top5Accuracy float64 `json:"top5_accuracy,omitempty"`
}

type SummaryPredictAccuracyInformations []SummaryPredictAccuracyInformation

func (SummaryPredictAccuracyInformation) Header() []string {
	extra := []string{
		"top1_accuracy",
		"top5_accuracy",
	}
	return append(SummaryBase{}.Header(), extra...)
}

func (s SummaryPredictAccuracyInformation) Row() []string {
	extra := []string{
		cast.ToString(s.Top1Accuracy),
		cast.ToString(s.Top5Accuracy),
	}
	return append(s.SummaryBase.Row(), extra...)
}

func (s SummaryPredictAccuracyInformation) key() string {
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

func (SummaryPredictAccuracyInformations) Header() []string {
	return SummaryPredictAccuracyInformation{}.Header()
}

func (s SummaryPredictAccuracyInformations) Rows() [][]string {
	rows := [][]string{}
	for _, e := range s {
		rows = append(rows, e.Row())
	}
	return rows
}

func (s SummaryPredictAccuracyInformations) Group() (SummaryPredictAccuracyInformations, error) {
	groups := map[string]SummaryPredictAccuracyInformations{}

	for _, v := range s {
		k := v.key()
		if _, ok := groups[k]; !ok {
			groups[k] = SummaryPredictAccuracyInformations{}
		}
		groups[k] = append(groups[k], v)
	}

	res := []SummaryPredictAccuracyInformation{}
	for _, v := range groups {
		if len(v) == 0 {
			log.Error("expecting more more than one input in SummaryPredictAccuracyInformations")
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

func (a ModelAccuracy) PredictAccuracyInformationSummary(e Evaluation) (*SummaryPredictAccuracyInformation, error) {
	return &SummaryPredictAccuracyInformation{
		SummaryBase:  e.summaryBase(),
		Top1Accuracy: a.Top1,
		Top5Accuracy: a.Top5,
	}, nil
}

func (as ModelAccuracies) PredictAccuracyInformationSummary(e Evaluation) ([]*SummaryPredictAccuracyInformation, error) {
	res := []*SummaryPredictAccuracyInformation{}
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

func (e Evaluation) PredictAccuracyInformationSummary(accCol *ModelAccuracyCollection) (*SummaryPredictAccuracyInformation, error) {
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

func (es Evaluations) PredictAccuracyInformationSummary(accCol *ModelAccuracyCollection) (SummaryPredictAccuracyInformations, error) {
	res := []SummaryPredictAccuracyInformation{}
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

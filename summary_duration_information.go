package evaluation

import (
	"errors"

	db "upper.io/db.v3"
)

type SummaryPredictDurationInformation struct {
	ID                  string
	ModelName           string
	ModelVersion        string
	FrameworkName       string
	FrameworkVersion    string
	FrameworkModel      string
	MachineArchitecture string
	UsingGPU            bool
	BatchSize           int
	HostName            string
	Durations           []uint64 // in nano seconds
}

func (p Performances) PredictDurationInformationSummary(e Evaluation) (*SummaryPredictDurationInformation, error) {
	spans := p.Spans().FilterByOperationName("predict")

	return &SummaryPredictDurationInformation{
		ID:                  e.ID,
		ModelName:           e.Model.Name,
		ModelVersion:        e.Model.Version,
		FrameworkName:       e.Framework.Name,
		FrameworkVersion:    e.Framework.Version,
		MachineArchitecture: e.MachineArchitecture,
		UsingGPU:            e.UsingGPU,
		BatchSize:           e.BatchSize,
		HostName:            e.Hostname,
		Durations:           spans.Duration(),
	}, nil
}

func (e Evaluation) PredictDurationInformationSummary(prefCol *PerformanceCollection) (*SummaryPredictDurationInformation, error) {
	perfs, err := prefCol.Find(db.Cond{"_id": e.PerformanceID})
	if err != nil {
		return nil, err
	}
	if len(perfs) != 1 {
		return nil, errors.New("expecting on performance output")
	}
	perf := perfs[0]
	return perf.PredictDurationInformationSummary(e)
}

func (es Evaluations) PredictDurationInformationSummary(prefCol *PerformanceCollection) ([]*SummaryPredictDurationInformation, error) {
	res := []*SummaryPredictDurationInformation{}
	for _, e := range es {
		s, err := e.PredictDurationInformationSummary(perfCol)
		if err != nil {
			log.WithError(err).Error("failed to get duration information summary")
			continue
		}
		res = append(res, s)
	}
	return nil, res
}

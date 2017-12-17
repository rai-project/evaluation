package evaluation

import (
	"errors"
	"time"

	db "upper.io/db.v3"
)

type SummaryPredictDurationInformation struct {
	SummaryBase
	MachineArchitecture string
	UsingGPU            bool
	BatchSize           int
	HostName            string
	Durations           []uint64 // in nano seconds
}

func (p Performance) PredictDurationInformationSummary(e Evaluation) (*SummaryPredictDurationInformation, error) {
	spans := p.Spans().FilterByOperationName("predict")

	return &SummaryPredictDurationInformation{
		SummaryBase: SummaryBase{
			ID:               e.ID,
			CreatedAt:        e.CreatedAt,
			UpdatedAt:        time.Now(),
			ModelName:        e.Model.Name,
			ModelVersion:     e.Model.Version,
			FrameworkName:    e.Framework.Name,
			FrameworkVersion: e.Framework.Version,
		},
		MachineArchitecture: e.MachineArchitecture,
		UsingGPU:            e.UsingGPU,
		BatchSize:           e.BatchSize,
		HostName:            e.Hostname,
		Durations:           spans.Duration(),
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

func (es Evaluations) PredictDurationInformationSummary(perfCol *PerformanceCollection) ([]*SummaryPredictDurationInformation, error) {
	res := []*SummaryPredictDurationInformation{}
	for _, e := range es {
		s, err := e.PredictDurationInformationSummary(perfCol)
		if err != nil {
			log.WithError(err).Error("failed to get duration information summary")
			continue
		}
		res = append(res, s)
	}
	return res, nil
}

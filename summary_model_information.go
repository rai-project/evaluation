package evaluation

import (
	"errors"
	"strings"
	"time"

	"github.com/spf13/cast"
	db "upper.io/db.v3"
)

//easyjson:json
type SummaryModelInformation struct {
	SummaryBase `json:",inline,omitempty"`
	Durations   []uint64 `json:"durations,omitempty"` // in nano seconds
}

type SummaryModelInformations []SummaryModelInformation

func (SummaryModelInformation) Header() []string {
	extra := []string{
		"durations (us)",
	}
	return append(SummaryBase{}.Header(), extra...)
}

func (s SummaryModelInformation) Row() []string {
	extra := []string{
		strings.Join(uint64SliceToStringSlice(s.Durations), ";"),
	}
	return append(s.SummaryBase.Row(), extra...)
}

func (SummaryModelInformations) Header() []string {
	return SummaryModelInformation{}.Header()
}

func (s SummaryModelInformations) Rows() [][]string {
	rows := [][]string{}
	for _, e := range s {
		rows = append(rows, e.Row())
	}
	return rows
}

func (p Performance) PredictDurationInformationSummary(e Evaluation) (*SummaryModelInformation, error) {
	spans := p.Spans().FilterByOperationName("c_predict")

	return &SummaryModelInformation{
		SummaryBase: e.summaryBase(),
		Durations:   spans.Duration(),
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

//easyjson:json
type SummaryModelInformationLatency struct {
	SummaryBase `json:",inline"`
	Durations   []float64 `json:"durations,omitempty"` // in nano seconds
	Duration    float64   `json:"duration,omitempty"`  // in nano seconds
	Latency     float64   `json:"latency,omitempty"`   // in nano seconds
	Throughput  float64   `json:"throughput,omitempty"`
}

//easyjson:json
type SummaryModelInformationLatencies []SummaryModelInformationLatency

func (SummaryModelInformationLatency) Header() []string {
	extra := []string{
		"durations",
		"duration (us)",
		"latency (us)",
		"throughput (input/s)",
	}
	return append(SummaryBase{}.Header(), extra...)
}

func (s SummaryModelInformationLatency) Row() []string {
	extra := []string{
		strings.Join(float64SliceToString(s.Durations), ";"),
		cast.ToString(s.Duration),
		cast.ToString(s.Latency),
		cast.ToString(s.Throughput),
	}
	return append(s.SummaryBase.Row(), extra...)
}

func (SummaryModelInformationLatencies) Header() []string {
	return SummaryModelInformationLatency{}.Header()
}

func (s SummaryModelInformationLatencies) Rows() [][]string {
	rows := [][]string{}
	for _, e := range s {
		rows = append(rows, e.Row())
	}
	return rows
}

func (info SummaryModelInformation) ThroughputLatencySummary() (SummaryModelInformationLatency, error) {
	var trimmedMeanFraction = DefaultTrimmedMeanFraction
	durations := toFloat64Slice(info.Durations)
	duration := TrimmedMean(durations, trimmedMeanFraction)
	return SummaryModelInformationLatency{
		SummaryBase: info.SummaryBase,
		Durations:   durations,
		Duration:    duration,
		Throughput:  float64(info.BatchSize) / duration,
		Latency:     duration / float64(info.BatchSize),
	}, nil
}

func (infos SummaryModelInformations) ThroughputLatencySummary() (SummaryModelInformationLatencies, error) {

	// MinIdx returns the index of the minimum value in the input slice. If several
	// entries have the maximum value, the first such index is returned. If the slice
	// is empty, MinIdx will panic.
	minIdx := func(s []float64) int {
		min := s[0]
		var ind int
		for i, v := range s {
			if v < min {
				min = v
				ind = i
			}
		}
		return ind
	}

	// Min returns the maximum value in the input slice. If the slice is empty, Min will panic.
	min := func(s []float64) float64 {
		return s[minIdx(s)]
	}

	var trimmedMeanFraction = DefaultTrimmedMeanFraction

	groups := map[string]SummaryModelInformations{}

	for _, info := range infos {
		k := info.key()
		if _, ok := groups[k]; !ok {
			groups[k] = SummaryModelInformations{}
		}
		groups[k] = append(groups[k], info)
	}

	res := []SummaryModelInformationLatency{}
	for _, v := range groups {
		if len(v) == 0 {
			log.Error("expecting more more than one input in SummaryModelInformationLatencies")
			continue
		}
		if len(v) == 1 {
			sum, err := v[0].ThroughputLatencySummary()
			if err != nil {
				log.WithError(err).Error("failed to get ThroughputLatencySummary")
				continue
			}
			res = append(res, sum)
			continue
		}

		durations := []float64{}
		for _, e := range v {
			if len(e.Durations) == 0 {
				continue
			}
			duration := TrimmedMean(toFloat64Slice(e.Durations), trimmedMeanFraction)
			if duration == 0 {
				continue
			}
			durations = append(durations, duration)
		}

		first := v[0]

		duration := min(durations)
		sum := SummaryModelInformationLatency{
			SummaryBase: first.SummaryBase,
			Durations:   durations,
			Duration:    duration,
			Throughput:  float64(first.BatchSize) * float64(time.Second/time.Microsecond) / duration,
			Latency:     duration / float64(first.BatchSize),
		}

		res = append(res, sum)
	}

	return res, nil
}

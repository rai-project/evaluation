package evaluation

import (
	"github.com/gonum/floats"
	"github.com/spf13/cast"
)

type SummaryThroughputLatency struct {
	SummaryBase `json:",inline"`
	Duration    float64 `json:"duration,omitempty"` // in nano seconds
	Latency     float64 `json:"latency,omitempty"`  // in nano seconds
	Throughput  float64 `json:"throughput,omitempty"`
}

type SummaryThroughputLatencies []SummaryThroughputLatency

func (SummaryThroughputLatency) Header() []string {
	extra := []string{
		"duration",
		"latency",
		"throughput",
	}
	return append(SummaryBase{}.Header(), extra...)
}

func (s SummaryThroughputLatency) Row() []string {
	extra := []string{
		cast.ToString(s.Duration),
		cast.ToString(s.Latency),
		cast.ToString(s.Throughput),
	}
	return append(s.SummaryBase.Row(), extra...)
}

func (SummaryThroughputLatencies) Header() []string {
	return SummaryThroughputLatency{}.Header()
}

func (s SummaryThroughputLatencies) Rows() [][]string {
	rows := [][]string{}
	for _, e := range s {
		rows = append(rows, e.Row())
	}
	return rows
}

func (info SummaryPredictDurationInformation) ThroughputLatencySummary() (SummaryThroughputLatency, error) {
	var trimmedMeanFraction = DefaultTrimmedMeanFraction
	duration := trimmedMean(toFloat64Slice(info.Durations), trimmedMeanFraction)
	return SummaryThroughputLatency{
		SummaryBase: info.SummaryBase,
		Duration:    duration,
		Throughput:  float64(info.BatchSize) / duration,
		Latency:     duration / float64(info.BatchSize),
	}, nil
}

func (infos SummaryPredictDurationInformations) ThroughputLatencySummary() (SummaryThroughputLatencies, error) {

	var trimmedMeanFraction = DefaultTrimmedMeanFraction

	groups := map[string]SummaryPredictDurationInformations{}

	for _, info := range infos {
		k := info.key()
		if _, ok := groups[k]; !ok {
			groups[k] = SummaryPredictDurationInformations{}
		}
		groups[k] = append(groups[k], info)
	}

	res := []SummaryThroughputLatency{}
	for _, v := range groups {
		if len(v) == 0 {
			log.Error("expecting more more than one input in SummaryThroughputLatencies")
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
			duration := trimmedMean(toFloat64Slice(e.Durations), trimmedMeanFraction)
			durations = append(durations, duration)
		}

		first := v[0]

		duration := floats.Min(durations)
		sum := SummaryThroughputLatency{
			SummaryBase: first.SummaryBase,
			Duration:    duration,
			Throughput:  float64(first.BatchSize) / duration,
			Latency:     duration / float64(first.BatchSize),
		}

		res = append(res, sum)
	}

	return res, nil
}

package evaluation

import (
	"strings"

	"github.com/gonum/floats"
	"github.com/spf13/cast"
)

type SummaryThroughputLatency struct {
	SummaryBase
	MachineArchitecture string
	UsingGPU            bool
	BatchSize           int
	HostName            string
	Duration            float64 // in nano seconds
	Latency             float64 // in nano seconds
	Throughput          float64
}

type SummaryThroughputLatencies []SummaryThroughputLatency

func (info SummaryPredictDurationInformation) ThroughputLatencySummary() (SummaryThroughputLatency, error) {
	var trimmedMeanFraction = DefaultTrimmedMeanFraction
	duration := trimmedMean(toFloat64Slice(info.Durations), trimmedMeanFraction)
	return SummaryThroughputLatency{
		SummaryBase:         info.SummaryBase,
		MachineArchitecture: info.MachineArchitecture,
		UsingGPU:            info.UsingGPU,
		BatchSize:           info.BatchSize,
		HostName:            info.HostName,
		Duration:            duration,
		Latency:             float64(info.BatchSize) / duration,
		Throughput:          duration / float64(info.BatchSize),
	}, nil
}

func (infos SummaryPredictDurationInformations) ThroughputLatencySummary() (SummaryThroughputLatencies, error) {

	var trimmedMeanFraction = DefaultTrimmedMeanFraction

	groups := map[string]SummaryPredictDurationInformations{}

	key := func(s SummaryPredictDurationInformation) string {
		return strings.Join(
			[]string{
				s.ModelName,
				s.ModelVersion,
				s.FrameworkName,
				s.FrameworkVersion,
				s.HostName,
				s.MachineArchitecture,
				cast.ToString(s.BatchSize),
				cast.ToString(s.UsingGPU),
			},
			",",
		)
	}

	for _, info := range infos {
		k := key(info)
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
			SummaryBase:         first.SummaryBase,
			MachineArchitecture: first.MachineArchitecture,
			UsingGPU:            first.UsingGPU,
			BatchSize:           first.BatchSize,
			HostName:            first.HostName,
			Duration:            duration,
			Latency:             float64(first.BatchSize) / duration,
			Throughput:          duration / float64(first.BatchSize),
		}

		res = append(res, sum)
	}

	return res, nil
}

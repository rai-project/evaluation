package evaluation

import (
	"strings"
	"time"

	"github.com/rai-project/evaluation/writer"
	"github.com/spf13/cast"
)

//easyjson:json
type SummaryModelLatencyInformation struct {
	SummaryBase `json:",inline"`
	Durations   []float64 `json:"durations,omitempty"`
	Duration    float64   `json:"duration,omitempty"`
	Latency     float64   `json:"latency,omitempty"`
	Throughput  float64   `json:"throughput,omitempty"`
}

//easyjson:json
type SummaryModelLatencyInformations []SummaryModelLatencyInformation

func (SummaryModelLatencyInformation) Header(opts ...writer.Option) []string {
	extra := []string{
		"durations",
		"duration (us)",
		"latency (us)",
		"throughput (input/s)",
	}
	return append(SummaryBase{}.Header(opts...), extra...)
}

func (s SummaryModelLatencyInformation) Row(opts ...writer.Option) []string {
	extra := []string{
		strings.Join(float64SliceToString(s.Durations), ";"),
		cast.ToString(s.Duration),
		cast.ToString(s.Latency),
		cast.ToString(s.Throughput),
	}
	return append(s.SummaryBase.Row(opts...), extra...)
}

func (SummaryModelLatencyInformations) Header(opts ...writer.Option) []string {
	return SummaryModelLatencyInformation{}.Header(opts...)
}

func (s SummaryModelLatencyInformations) Rows(opts ...writer.Option) [][]string {
	rows := [][]string{}
	for _, e := range s {
		rows = append(rows, e.Row(opts...))
	}
	return rows
}

func (info SummaryModelInformation) ThroughputLatencySummary() (SummaryModelLatencyInformation, error) {
	var trimmedMeanFraction = DefaultTrimmedMeanFraction
	durations := toFloat64Slice(info.Durations)
	duration := TrimmedMean(durations, trimmedMeanFraction)
	return SummaryModelLatencyInformation{
		SummaryBase: info.SummaryBase,
		Durations:   durations,
		Duration:    duration,
		Throughput:  float64(info.BatchSize) / duration,
		Latency:     duration / float64(info.BatchSize),
	}, nil
}

func (infos SummaryModelInformations) ThroughputLatencySummary() (SummaryModelLatencyInformations, error) {

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

	res := []SummaryModelLatencyInformation{}
	for _, v := range groups {
		if len(v) == 0 {
			log.Error("expecting more more than one input in SummaryModelLatencyInformations")
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
		sum := SummaryModelLatencyInformation{
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

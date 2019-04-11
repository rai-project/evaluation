package plotting

import (
	"github.com/fatih/structs"
	"github.com/spf13/cast"
)

type batchPlot struct {
	Durations []batchDurationSummary
	Options   *Options
}

type batchDurationSummary struct {
	BatchSize int
	Duration  uint64
}

func NewBatchPlot(os ...OptionModifier) (*batchPlot, error) {
	os = append(os, Option.BatchSize(0))
	opts := NewOptions(os...)

	traceFilePaths, err := modelTraceFiles(opts)
	if err != nil {
		return nil, err
	}

	durations := []batchDurationSummary{}
	for _, path := range traceFilePaths {
		tr, err := ReadTraceFile(path)
		if err != nil {
			if opts.ignoreReadErrors {
				continue
			}
			return nil, err
		}
		cPredictSpan, err := findCPredict(tr)
		if err != nil {
			return nil, err
		}
		batchSize, err := findBatchSize(tr)
		if err != nil {
			return nil, err
		}
		durations = append(durations,
			batchDurationSummary{
				BatchSize: batchSize,
				Duration:  cPredictSpan.Duration,
			},
		)
	}
	return &batchPlot{
		Durations: durations,
		Options:   opts,
	}, nil
}

func (o batchPlot) Header() []string {
	return batchDurationSummary{}.Header()
}

func (o batchPlot) Rows() [][]string {
	res := make([][]string, len(o.Durations))
	for ii, dur := range o.Durations {
		res[ii] = dur.Row()
	}
	return res
}

func (o batchDurationSummary) Header() []string {
	res := []string{}
	s := structs.New(&o)
	for _, field := range s.Fields() {
		res = append(res, field.Name())
	}
	return res
}

func (o batchDurationSummary) Row() []string {
	res := []string{}
	s := structs.New(&o)
	for _, field := range s.Fields() {
		res = append(res, cast.ToString(field.Value()))
	}
	return res
}

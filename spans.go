package evaluation

import (
	model "github.com/uber/jaeger/model/json"
)

type Spans []model.Span

func (spns Spans) FilterByOperationName(op string) Spans {
	res := []model.Span{}
	for _, s := range spns {
		if s.OperationName == op {
			res = append(res, s)
		}
	}
	return res
}

func (spns Spans) Duration() []uint64 {
	res := make([]uint64, len(spns))
	for ii, s := range spns {
		res[ii] = uint64(s.Duration)
	}
	return res
}

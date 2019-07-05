package evaluation

import (
	"strings"

	model "github.com/uber/jaeger/model/json"
)

type Spans []model.Span

func (spns Spans) FilterByOperationName(op string) Spans {
	res := []model.Span{}
	op = strings.ToLower(op)
	for _, s := range spns {
		if strings.ToLower(s.OperationName) == op {
			res = append(res, s)
		}
	}
	return res
}

func (spns Spans) FilterByOperationNameAndEvalTraceLevel(op string, lvl string) Spans {
	res := []model.Span{}
	op = strings.ToLower(op)
	for _, s := range spns {
		traceLevel, err := getTagValueAsString(s, "evaluation_trace_level")
		if err != nil || traceLevel == "" {
			continue
		}
		if strings.ToLower(s.OperationName) == op && traceLevel == lvl {
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

package evaluation

import model "github.com/uber/jaeger/model/json"

type Performances []Performance

func (ps Performances) Spans() Spans {
	res := []model.Span{}
	for _, p := range ps {
		res = append(res, p.Spans()...)
	}
	return res
}

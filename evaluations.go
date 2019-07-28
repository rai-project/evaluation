package evaluation

import (
	"errors"

	"github.com/k0kubun/pp"
	model "github.com/uber/jaeger/model/json"
	"upper.io/db.v3"
)

type Evaluations []Evaluation

func (es Evaluations) GetSpansFromPerformanceCollection(perfCol *PerformanceCollection) (Spans, error) {
	spans := []model.Span{}
	for _, e := range es {
		foundPerfs, err := perfCol.Find(db.Cond{"_id": e.PerformanceID})
		if err != nil {
			return nil, err
		}
		if len(foundPerfs) != 1 {
			return nil, errors.New("expecting on performance output")
		}
		perf := foundPerfs[0]
		pp.Println(len(perf.Spans()))
		spans = append(spans, perf.Spans()...)
	}
	return spans, nil
}

func (es Evaluations) GroupByBatchSize() map[int]Evaluations {
	ret := make(map[int]Evaluations)
	for _, e := range es {
		_, ok := ret[e.BatchSize]
		if !ok {
			ret[e.BatchSize] = Evaluations{e}
		} else {
			ret[e.BatchSize] = append(ret[e.BatchSize], e)
		}
	}
	return ret
}

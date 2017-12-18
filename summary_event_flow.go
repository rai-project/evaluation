package evaluation

import (
	"errors"

	"github.com/rai-project/evaluation/eventflow"
	db "upper.io/db.v3"
)

type SummaryEventFlow struct {
	SummaryBase `json:",inline"`
	EventFlow   eventflow.Event `json:"event_flow,omitempty"`
}

type SummaryEventFlows []SummaryEventFlow

func spansToEventFlow(spans Spans) eventflow.Events {
	return eventflow.SpansToEvenFlow(spans)
}

func (p Performance) EventFlowSummary(e Evaluation) (*SummaryEventFlow, error) {
	flow := spansToEventFlow(p.Spans())
	return &SummaryEventFlow{
		SummaryBase: e.summaryBase(),
		EventFlow:   flow,
	}, nil
}

func (e Evaluation) EventFlowSummary(perfCol *PerformanceCollection) (*SummaryEventFlow, error) {
	perfs, err := perfCol.Find(db.Cond{"_id": e.PerformanceID})
	if err != nil {
		return nil, err
	}
	if len(perfs) != 1 {
		return nil, errors.New("expecting on performance output")
	}
	perf := perfs[0]
	return perf.EventFlowSummary(e)
}

func (es Evaluations) EventFlowSummary(perfCol *PerformanceCollection) (SummaryEventFlows, error) {
	res := []SummaryLayerInformation{}
	for _, e := range es {
		s, err := e.EventFlowSummary(perfCol)
		if err != nil {
			log.WithError(err).Error("failed to get layer information summary")
			continue
		}
		if s == nil {
			continue
		}
		res = append(res, *s)
	}
	return res, nil
}

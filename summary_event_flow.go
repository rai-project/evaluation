package evaluation

import (
	"github.com/rai-project/evaluation/eventflow"
)

func spansToEventFlow(spans Spans) eventflow.Events {
	return eventflow.SpansToEvenFlow(spans)
}

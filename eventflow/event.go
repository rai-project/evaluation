package eventflow

import (
	"encoding/json"
	"time"

	"github.com/rai-project/evaluation/writer"
	"github.com/spf13/cast"
	m "github.com/uber/jaeger/model"
	model "github.com/uber/jaeger/model/json"
)

//easyjson:json
type Event struct {
	ID        string            `json:"EVENT_ID"`
	ParentID  string            `json:"PARENT_ID"`
	Name      string            `json:"EVENT_NAME"`
	MetaData  map[string]string `json:"META,omitempty"`
	TimeStamp time.Time         `json:"TS,omitempty"`
	Duration  uint64            `json:"ELAPSED_MS,omitempty"`
}

type Events []Event

func (Event) Header(opts ...writer.Option) []string {
	return []string{
		"id",
		"parent_id",
		"name",
		"metadata",
		"timestamp",
		"duration (us)",
	}
}

func (e Event) Row(opts ...writer.Option) []string {
	metadata, err := json.Marshal(e.MetaData)
	if err != nil {
		metadata = []byte{}
	}
	return []string{
		e.ID,
		e.ParentID,
		e.Name,
		string(metadata),
		e.TimeStamp.String(),
		cast.ToString(e.Duration),
	}
}

func (Events) Header(opts ...writer.Option) []string {
	return Event{}.Header(opts...)
}

func (s Events) Rows(opts ...writer.Option) [][]string {
	rows := [][]string{}
	for _, e := range s {
		rows = append(rows, e.Row(opts...))
	}
	return rows
}

func tagsOf(span model.Span) map[string]string {
	res := map[string]string{}
	for _, lg := range span.Logs {
		for _, fld := range lg.Fields {
			res[fld.Key] = cast.ToString(fld.Value)
		}
	}
	for _, tag := range span.Tags {
		res[tag.Key] = cast.ToString(tag.Value)
	}
	return res
}

func parentOf(span model.Span) model.SpanID {
	if span.ParentSpanID != "" {
		return span.ParentSpanID
	}
	for _, ref := range span.References {
		if ref.RefType == model.ChildOf {
			return ref.SpanID
		}
	}
	return model.SpanID("")
}

func toTime(t uint64) time.Time {
	return m.EpochMicrosecondsAsTime(t)
}

func toDuration(d uint64) time.Duration {
	return m.MicrosecondsAsDuration(d)
}

func spanToEvent(span model.Span) Event {
	return Event{
		ID:        string(span.SpanID),
		ParentID:  string(parentOf(span)),
		Name:      span.OperationName,
		MetaData:  tagsOf(span),
		TimeStamp: toTime(span.StartTime),
		Duration:  span.Duration,
	}
}

func SpansToEvenFlow(spans []model.Span) Events {
	events := make([]Event, len(spans))
	for ii, span := range spans {
		events[ii] = spanToEvent(span)
	}
	return events
}

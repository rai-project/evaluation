package evaluation

import (
	"time"

	"github.com/rai-project/tracer"
	"gopkg.in/mgo.v2/bson"

	model "github.com/uber/jaeger/model/json"
)

type TraceInformation struct {
	json   string            `json:"-"`
	Traces []model.Trace     `json:"data"`
	Total  int               `json:"total"`
	Limit  int               `json:"limit"`
	Offset int               `json:"offset"`
	Errors []structuredError `json:"errors"`
}

type structuredError struct {
	Code    int           `json:"code,omitempty"`
	Msg     string        `json:"msg"`
	TraceID model.TraceID `json:"traceID,omitempty"`
}

type Performance struct {
	ID         bson.ObjectId `json:"id" bson:"_id"`
	CreatedAt  time.Time     `json:"created_at"  bson:"created_at"`
	Trace      TraceInformation
	TraceLevel tracer.Level
}

func (Performance) TableName() string {
	return "performance"
}

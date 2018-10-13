package evaluation

import (
	"time"

	"github.com/rai-project/database"
	"github.com/rai-project/database/mongodb"
	"github.com/rai-project/tracer"
	"gopkg.in/mgo.v2/bson"

	model "github.com/uber/jaeger/model/json"
)

//easyjson:json
type TraceInformation struct {
	Traces []model.Trace     `json:"data"`
	Total  int               `json:"total"`
	Limit  int               `json:"limit"`
	Offset int               `json:"offset"`
	Errors []structuredError `json:"errors"`
}

func (info TraceInformation) Spans() Spans {
	res := []model.Span{}
	for _, tr := range info.Traces {
		res = append(res, tr.Spans...)
	}
	return Spans(res)
}

//easyjson:json
type structuredError struct {
	Code    int           `json:"code,omitempty"`
	Msg     string        `json:"msg"`
	TraceID model.TraceID `json:"traceID,omitempty"`
}

//easyjson:json
type Performance struct {
	ID         bson.ObjectId `json:"id" bson:"_id"`
	CreatedAt  time.Time     `json:"created_at"  bson:"created_at"`
	Trace      TraceInformation
	TraceLevel tracer.Level
}

func (Performance) TableName() string {
	return "performance"
}

func (p Performance) Spans() Spans {
	return p.Trace.Spans()
}

type PerformanceCollection struct {
	*mongodb.MongoTable
}

func NewPerformanceCollection(db database.Database) (*PerformanceCollection, error) {
	tbl, err := mongodb.NewTable(db, Performance{}.TableName())
	if err != nil {
		return nil, err
	}
	tbl.Create(nil)

	return &PerformanceCollection{
		MongoTable: tbl.(*mongodb.MongoTable),
	}, nil
}

func (c *PerformanceCollection) Find(as ...interface{}) ([]Performance, error) {
	pref := []Performance{}

	collection := c.Session.Collection(c.Name())

	err := collection.Find(as...).All(&pref)
	if err != nil {
		return nil, err
	}
	return pref, nil
}

func (m *PerformanceCollection) Close() error {
	return nil
}

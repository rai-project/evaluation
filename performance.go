package evaluation

import (
	"time"

	"github.com/rai-project/database"
	"github.com/rai-project/database/mongodb"
	"github.com/rai-project/tracer"
	"gopkg.in/mgo.v2/bson"

	model "github.com/uber/jaeger/model/json"
)

type TraceInformation struct {
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

func (m *PerformanceCollection) Close() error {
	return nil
}

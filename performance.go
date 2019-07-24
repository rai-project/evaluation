package evaluation

import (
	"encoding/base64"
	json "encoding/json"
	"time"

	"github.com/pkg/errors"

	"github.com/golang/snappy"
	"github.com/rai-project/database"
	"github.com/rai-project/database/mongodb"
	"github.com/rai-project/tracer"
	"gopkg.in/mgo.v2/bson"

	model "github.com/uber/jaeger/model/json"
)

//easyjson:json
type TraceInformation struct {
	Traces []model.Trace     `bson:"data,omitempty" json:"data,omitempty"`
	Total  int               `bson:"total,omitempty" json:"total,omitempty"`
	Limit  int               `bson:"limit,omitempty" json:"limit,omitempty"`
	Offset int               `bson:"offset,omitempty" json:"offset,omitempty"`
	Errors []structuredError `bson:"errors,omitempty" json:"errors,omitempty"`
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
	Code    int           `json:"code,omitempty" bson:"code,omitempty"`
	Msg     string        `json:"msg,omitempty" bson:"msg,omitempty"`
	TraceID model.TraceID `json:"traceID,omitempty" bson:"traceID,omitempty"`
}

//easyjson:json
type Performance struct {
	ID              bson.ObjectId     `json:"id,omitempty"  bson:"_id,omitempty"`
	CreatedAt       time.Time         `json:"created_at" bson:"created_at,omitempty"`
	TraceCompressed []byte            `json:"trace_compressed" bson:"trace,omitempty"`
	Trace           *TraceInformation `json:"trace" bson:"-"`
	TraceLevel      tracer.Level      `json:"trace_level" bson:"trace_level,omitempty"`
}

func (Performance) TableName() string {
	return "performance"
}

func (p Performance) Spans() Spans {
	return p.Trace.Spans()
}

func (p *Performance) UncompressTrace() error {
	if p.TraceCompressed == nil || len(p.TraceCompressed) == 0 {
		return nil
	}
	bts, err := base64.StdEncoding.DecodeString(string(p.TraceCompressed))
	if err != nil {
		return errors.Wrap(err, "cannot decode compressed trace")
	}
	out, err := snappy.Decode(nil, bts)
	if err != nil {
		return errors.Wrap(err, "cannot snappy uncompressed trace")
	}

	traceInfo := &TraceInformation{}

	err = json.Unmarshal(out, traceInfo)
	if err != nil {
		return errors.Wrap(err, "cannot unmarshal trace")
	}
	p.Trace = traceInfo

	return nil
}

func (p *Performance) CompressTrace() error {
	bts, err := json.Marshal(p.Trace)
	if err != nil {
		return errors.Wrap(err, "cannot marshal trace")
	}
	out := snappy.Encode(nil, bts)
	p.TraceCompressed = []byte(base64.StdEncoding.EncodeToString(out))
	p.Trace = nil
	return nil
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
	perfs := []Performance{}

	collection := c.Session.Collection(c.Name())

	err := collection.Find(as...).All(&perfs)
	if err != nil {
		return nil, err
	}
	for ii := range perfs {
		err = (&perfs[ii]).UncompressTrace()
		if err != nil {
			return nil, errors.New("unable to uncompress trace")
		}
	}
	return perfs, nil
}

func (m *PerformanceCollection) Close() error {
	return nil
}

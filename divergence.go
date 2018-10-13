package evaluation

import (
	"time"

	"github.com/rai-project/database"
	"github.com/rai-project/database/mongodb"
	"github.com/rai-project/dlframework"
	"gopkg.in/mgo.v2/bson"
)

//easyjson:json
type Divergence struct {
	ID                           bson.ObjectId `json:"id" bson:"_id"`
	CreatedAt                    time.Time     `json:"created_at"  bson:"created_at"`
	Method                       string
	Value                        float64
	SourcePredictionID           bson.ObjectId
	TargetPredictionID           bson.ObjectId
	SourceInputPredictionInputID string
	TargetInputPredictionInputID string
	SourceFeatures               dlframework.Features
	TargetFeatures               dlframework.Features
}

func (Divergence) TableName() string {
	return "divergence"
}

type DivergenceCollection struct {
	*mongodb.MongoTable
}

func NewDivergenceCollection(db database.Database) (*DivergenceCollection, error) {
	tbl, err := mongodb.NewTable(db, Divergence{}.TableName())
	if err != nil {
		return nil, err
	}
	tbl.Create(nil)

	return &DivergenceCollection{
		MongoTable: tbl.(*mongodb.MongoTable),
	}, nil
}

func (m *DivergenceCollection) Close() error {
	return nil
}

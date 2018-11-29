package evaluation

import (
	"time"

	"github.com/rai-project/database"
	"github.com/rai-project/database/mongodb"
	"github.com/rai-project/dlframework"
	"gopkg.in/mgo.v2/bson"
)

//easybson:json
type Divergence struct {
	ID                           bson.ObjectId        `bson:"id,omitempty" json:"id,omitempty"`
	CreatedAt                    time.Time            `bson:"created_at,omitempty" bson:"created_at" json:"created_at,omitempty"`
	Method                       string               `bson:"method,omitempty" json:"method,omitempty"`
	Value                        float64              `bson:"value,omitempty" json:"value,omitempty"`
	SourcePredictionID           bson.ObjectId        `bson:"source_prediction_id,omitempty" json:"source_prediction_id,omitempty"`
	TargetPredictionID           bson.ObjectId        `bson:"target_prediction_id,omitempty" json:"target_prediction_id,omitempty"`
	SourceInputPredictionInputID string               `bson:"source_input_prediction_input_id,omitempty" json:"source_input_prediction_input_id,omitempty"`
	TargetInputPredictionInputID string               `bson:"target_input_prediction_input_id,omitempty" json:"target_input_prediction_input_id,omitempty"`
	SourceFeatures               dlframework.Features `bson:"source_features,omitempty" json:"source_features,omitempty"`
	TargetFeatures               dlframework.Features `bson:"target_features,omitempty" json:"target_features,omitempty"`
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

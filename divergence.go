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
	ID                           bson.ObjectId        `json:"id,omitempty" bson:"_id"`
	CreatedAt                    time.Time            `json:"created_at,omitempty" bson:"created_at"`
	Method                       string               `json:"method,omitempty"`
	Value                        float64              `json:"value,omitempty"`
	SourcePredictionID           bson.ObjectId        `json:"source_prediction_id,omitempty"`
	TargetPredictionID           bson.ObjectId        `json:"target_prediction_id,omitempty"`
	SourceInputPredictionInputID string               `json:"source_input_prediction_input_id,omitempty"`
	TargetInputPredictionInputID string               `json:"target_input_prediction_input_id,omitempty"`
	SourceFeatures               dlframework.Features `json:"source_features,omitempty"`
	TargetFeatures               dlframework.Features `json:"target_features,omitempty"`
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

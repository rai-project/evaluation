package evaluation

import (
	"time"

	"github.com/rai-project/dlframework"
	"gopkg.in/mgo.v2/bson"
)

type InputPrediction struct {
	ID            bson.ObjectId `json:"id" bson:"_id"`
	CreatedAt     time.Time     `json:"created_at"  bson:"created_at"`
	InputID       string
	ExpectedLabel string
	Predictions   dlframework.Features
}

func (InputPrediction) TableName() string {
	return "input_prediction"
}

type ModelAccuracy struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	CreatedAt time.Time     `json:"created_at"  bson:"created_at"`
	Top1      float64
	Top5      float64
}

func (ModelAccuracy) TableName() string {
	return "model_accuracy"
}

package evaluation

import (
	"time"

	"github.com/rai-project/dlframework"
	"gopkg.in/mgo.v2/bson"
)

type Evaluation struct {
	ID                 bson.ObjectId `json:"id" bson:"_id"`
	CreatedAt          time.Time     `json:"created_at"  bson:"created_at"`
	Framework          dlframework.FrameworkManifest
	Model              dlframework.ModelManifest
	DatasetCategory    string
	DatasetName        string
	ModelAccuracyID    bson.ObjectId
	InputPredictionIDs []bson.ObjectId
	PerformanceID      bson.ObjectId
	Public             bool
}

package evaluation

import (
	"time"

	"github.com/rai-project/dlframework"
	"gopkg.in/mgo.v2/bson"
)

type Evaluation struct {
	ID                  bson.ObjectId `json:"id" bson:"_id"`
	CreatedAt           time.Time     `json:"created_at"  bson:"created_at"`
	Framework           dlframework.FrameworkManifest
	Model               dlframework.ModelManifest
	DatasetCategory     string
	DatasetName         string
	MachineArchitecture string
	UsingGPU            bool
	BatchSize           int
	Hostname            string
	TraceLevel          string
	ModelAccuracyID     bson.ObjectId
	InputPredictionIDs  []bson.ObjectId
	PerformanceID       bson.ObjectId
	Public              bool
	Metadata            map[string]string
}

func (Evaluation) TableName() string {
	return "evaluation"
}

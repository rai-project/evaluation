package evaluation

import (
	"time"

	"github.com/rai-project/database"
	"github.com/rai-project/database/mongodb"
	"github.com/rai-project/dlframework"
	"gopkg.in/mgo.v2/bson"
	"upper.io/db.v3"
)

//easyjson:json
type Evaluation struct {
	ID                  bson.ObjectId                 `json:"id,omitempty" bson:"_id"`
	UserID              string                        `json:"user_id,omitempty"`
	RunID               int                           `json:"run_id,omitempty"`
	CreatedAt           time.Time                     `json:"created_at,omitempty" bson:"created_at"`
	Framework           dlframework.FrameworkManifest `json:"framework,omitempty"`
	Model               dlframework.ModelManifest     `json:"model,omitempty"`
	DatasetCategory     string                        `json:"dataset_category,omitempty"`
	DatasetName         string                        `json:"dataset_name,omitempty"`
	MachineArchitecture string                        `json:"machine_architecture,omitempty"`
	UsingGPU            bool                          `json:"using_gpu,omitempty"`
	BatchSize           int                           `json:"batch_size,omitempty"`
	Hostname            string                        `json:"hostname,omitempty"`
	TraceLevel          string                        `json:"trace_level,omitempty"`
	ModelAccuracyID     bson.ObjectId                 `json:"model_accuracy_id,omitempty"`
	InputPredictionIDs  []bson.ObjectId               `json:"input_prediction_i_ds,omitempty"`
	PerformanceID       bson.ObjectId                 `json:"performance_id,omitempty"`
	Public              bool                          `json:"public,omitempty"`
	Metadata            map[string]string             `json:"metadata,omitempty"`
}

func (Evaluation) TableName() string {
	return "evaluation"
}

type EvaluationCollection struct {
	*mongodb.MongoTable
}

func NewEvaluationCollection(db database.Database) (*EvaluationCollection, error) {
	tbl, err := mongodb.NewTable(db, Evaluation{}.TableName())
	if err != nil {
		return nil, err
	}
	tbl.Create(nil)

	return &EvaluationCollection{
		MongoTable: tbl.(*mongodb.MongoTable),
	}, nil
}

func (c *EvaluationCollection) Find(as ...interface{}) ([]Evaluation, error) {
	evals := []Evaluation{}

	collection := c.Session.Collection(c.Name())

	err := collection.Find(as...).All(&evals)
	if err != nil {
		return nil, err
	}
	return evals, nil
}

func (c *EvaluationCollection) FindByModel(model dlframework.ModelManifest) ([]Evaluation, error) {
	return c.Find(
		db.Cond{
			"model.name":    model.Name,
			"model.version": model.Version,
		},
	)
}

func (m *EvaluationCollection) Close() error {
	return nil
}

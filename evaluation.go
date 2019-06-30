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
	ID                  bson.ObjectId                 `json:"id,omitempty" bson:"_id,omitempty"`
	UserID              string                        `json:"user_id,omitempty" bson:"user_id,omitempty"`
	RunID               string                        `json:"run_id,omitempty"  bson:"run_id,omitempty"`
	CreatedAt           time.Time                     `json:"created_at,omitempty" bson:"created_at,omitempty"`
	Framework           dlframework.FrameworkManifest `json:"framework,omitempty"  bson:"framework,omitempty"`
	Model               dlframework.ModelManifest     `json:"model,omitempty"  bson:"model,omitempty"`
	DatasetCategory     string                        `json:"dataset_category,omitempty" bson:"dataset_category"`
	DatasetName         string                        `json:"dataset_name,omitempty" bson:"dataset_name,omitempty"`
	MachineArchitecture string                        `json:"machine_architecture,omitempty" bson:"machine_architecture,omitempty"`
	UsingGPU            bool                          `json:"using_gpu,omitempty" bson:"using_gpu,omitempty"`
	BatchSize           int                           `json:"batch_size,omitempty" bson:"batch_size,omitempty"`
	GPUMetrics          string                        `json:"gpu_metrics,omitempty" bson:"gpu_metrics,omitempty"`
	Hostname            string                        `json:"hostname,omitempty" bson:"hostname,omitempty"`
	TraceLevel          string                        `json:"trace_level,omitempty" bson:"trace_level,omitempty"`
	ModelAccuracyID     bson.ObjectId                 `json:"model_accuracy_id,omitempty" bson:"model_accuracy_id,omitempty"`
	InputPredictionIDs  []bson.ObjectId               `json:"input_prediction_ids,omitempty" bson:"input_prediction_ids,omitempty"`
	PerformanceID       bson.ObjectId                 `json:"performance_id,omitempty" bson:"performance_id,omitempty"`
	Public              bool                          `json:"public,omitempty" bson:"public,omitempty"`
	Metadata            map[string]string             `json:"metadata,omitempty" bson:"metadata,omitempty"`
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

func (c *EvaluationCollection) FindByUserID(UserID string) ([]Evaluation, error) {
	return c.Find(
		db.Cond{
			"userid": UserID,
		},
	)
}

func (m *EvaluationCollection) Close() error {
	return nil
}

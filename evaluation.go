package evaluation

import (
	"time"

	"github.com/rai-project/database"
	"github.com/rai-project/database/mongodb"
	"github.com/rai-project/dlframework"
	"gopkg.in/mgo.v2/bson"
	"upper.io/db.v3"
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

func (c EvaluationCollection) Find(as ...interface{}) ([]Evaluation, error) {
	evals := []Evaluation{}

	collection := c.Session.Collection(c.Name())

	err := collection.Find(as...).All(&evals)
	if err != nil {
		return nil, err
	}
	return evals, nil
}

func (c EvaluationCollection) FindByModel(model dlframework.ModelManifest) ([]Evaluation, error) {
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

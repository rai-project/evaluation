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
type ModelAccuracy struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	CreatedAt time.Time     `json:"created_at"  bson:"created_at"`
	Top1      float64
	Top5      float64
}

func (ModelAccuracy) TableName() string {
	return "model_accuracy"
}

type ModelAccuracyCollection struct {
	*mongodb.MongoTable
}

func NewModelAccuracyCollection(db database.Database) (*ModelAccuracyCollection, error) {
	tbl, err := mongodb.NewTable(db, ModelAccuracy{}.TableName())
	if err != nil {
		return nil, err
	}
	tbl.Create(nil)

	return &ModelAccuracyCollection{
		MongoTable: tbl.(*mongodb.MongoTable),
	}, nil
}

func (c *ModelAccuracyCollection) Find(as ...interface{}) ([]ModelAccuracy, error) {
	accs := []ModelAccuracy{}

	collection := c.Session.Collection(c.Name())

	err := collection.Find(as...).All(&accs)
	if err != nil {
		return nil, err
	}
	return accs, nil
}

func (c *ModelAccuracyCollection) FindByModel(model dlframework.ModelManifest) ([]ModelAccuracy, error) {
	return c.Find(
		db.Cond{
			"model.name":    model.Name,
			"model.version": model.Version,
		},
	)
}

func (m *ModelAccuracyCollection) Close() error {
	return nil
}

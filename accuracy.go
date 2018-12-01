package evaluation

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"
	"upper.io/db.v3"

	"github.com/rai-project/database"
	"github.com/rai-project/database/mongodb"
	"github.com/rai-project/dlframework"
)

//easybson:json
type ModelAccuracy struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	CreatedAt time.Time     `bson:"created_at,omitempty" json:"created_at,omitempty"`
	Top1      float64       `bson:"top_1,omitempty" json:"top_1,omitempty"`
	Top5      float64       `bson:"top_5,omitempty" json:"top_5,omitempty"`
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
		fmt.Print("here")
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

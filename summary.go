package evaluation

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type SummaryBase struct {
	ID               bson.ObjectId `json:"id" bson:"_id"`
	CreatedAt        time.Time     `json:"created_at"  bson:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"  bson:"updated_at"`
	ModelName        string
	ModelVersion     string
	FrameworkName    string
	FrameworkVersion string
	FrameworkModel   string
}

func (e Evaluation) summaryBase() SummaryBase {
	return SummaryBase{
		ID:               e.ID,
		CreatedAt:        e.CreatedAt,
		UpdatedAt:        time.Now(),
		ModelName:        e.Model.Name,
		ModelVersion:     e.Model.Version,
		FrameworkName:    e.Framework.Name,
		FrameworkVersion: e.Framework.Version,
	}
}

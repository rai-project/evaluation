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
}

func (SummaryBase) Header() []string {
	return []string{
		"id",
		"created_at",
		"updated_at",
		"model_name",
		"model_version",
		"framework_name",
		"framework_version",
	}
}

func (s SummaryBase) Row() []string {
	return []string{
		s.ID.String(),
		s.CreatedAt.String(),
		s.UpdatedAt.String(),
		s.ModelName,
		s.ModelVersion,
		s.FrameworkName,
		s.FrameworkVersion,
	}
}

func (s SummaryBase) FrameworkModel() string {
	return s.FrameworkName + "::" + s.FrameworkVersion + "/" + s.ModelName + "::" + s.ModelVersion
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

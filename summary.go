package evaluation

import (
	"fmt"
	"strings"
	"time"

	"github.com/rai-project/evaluation/writer"
	"github.com/spf13/cast"
	"gopkg.in/mgo.v2/bson"
)

//easyjson:json
type SummaryBase struct {
	ID                  bson.ObjectId `json:"id" bson:"_id"`
	CreatedAt           time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt           time.Time     `json:"updated_at" bson:"updated_at"`
	ModelName           string        `json:"model_name,omitempty"`
	ModelVersion        string        `json:"model_version,omitempty"`
	FrameworkName       string        `json:"framework_name,omitempty"`
	FrameworkVersion    string        `json:"framework_version,omitempty"`
	MachineArchitecture string        `json:"machine_architecture,omitempty"`
	UsingGPU            bool          `json:"using_gpu,omitempty"`
	BatchSize           int           `json:"batch_size,omitempty"`
	HostName            string        `json:"host_name,omitempty"`
}

func (SummaryBase) Header(opts ...writer.Option) []string {
	return []string{
		"id",
		"created_at",
		"updated_at",
		"model_name",
		"model_version",
		"framework_name",
		"framework_version",
		"machine_architecture",
		"using_gpu",
		"batch_size",
		"hostname",
	}
}

func (s SummaryBase) Row(opts ...writer.Option) []string {
	return []string{
		fmt.Sprintf(`%x`, string(s.ID)),
		s.CreatedAt.String(),
		s.UpdatedAt.String(),
		s.ModelName,
		s.ModelVersion,
		s.FrameworkName,
		s.FrameworkVersion,
		s.MachineArchitecture,
		cast.ToString(s.UsingGPU),
		cast.ToString(s.BatchSize),
		s.HostName,
	}
}

func (s SummaryBase) FrameworkModel() string {
	return s.FrameworkName + "::" + s.FrameworkVersion + "/" + s.ModelName + "::" + s.ModelVersion
}

func (s SummaryBase) key() string {
	return strings.Join(
		[]string{
			s.ModelName,
			s.ModelVersion,
			s.FrameworkName,
			s.FrameworkVersion,
			s.HostName,
			s.MachineArchitecture,
			cast.ToString(s.BatchSize),
			cast.ToString(s.UsingGPU),
		},
		",",
	)
}

func (e Evaluation) summaryBase() SummaryBase {
	return SummaryBase{
		ID:                  e.ID,
		CreatedAt:           e.CreatedAt,
		UpdatedAt:           time.Now(),
		ModelName:           e.Model.Name,
		ModelVersion:        e.Model.Version,
		FrameworkName:       e.Framework.Name,
		FrameworkVersion:    e.Framework.Version,
		MachineArchitecture: e.MachineArchitecture,
		UsingGPU:            e.UsingGPU,
		BatchSize:           e.BatchSize,
		HostName:            e.Hostname,
	}
}

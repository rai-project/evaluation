package evaluation

import (
	"fmt"
	"strings"
	"time"

	"github.com/rai-project/evaluation/writer"
	"github.com/rai-project/machine"
	nvidiasmi "github.com/rai-project/nvidia-smi"
	"github.com/spf13/cast"
	"gopkg.in/mgo.v2/bson"
)

var (
	cntkLogMessageShown = false
)

//easyjson:json
type SummaryBase struct {
	ID                       bson.ObjectId    `json:"id" bson:"_id"`
	CreatedAt                time.Time        `json:"created_at" bson:"created_at"`
	UpdatedAt                time.Time        `json:"updated_at" bson:"updated_at"`
	ModelName                string           `json:"model_name,omitempty"`
	ModelVersion             string           `json:"model_version,omitempty"`
	FrameworkName            string           `json:"framework_name,omitempty"`
	FrameworkVersion         string           `json:"framework_version,omitempty"`
	MachineArchitecture      string           `json:"machine_architecture,omitempty"`
	UsingGPU                 bool             `json:"using_gpu,omitempty"`
	BatchSize                int              `json:"batch_size,omitempty"`
	HostName                 string           `json:"host_name,omitempty"`
	TraceLevel               string           `json:"trace_level,omitempty"`
	MachineInformation       *machine.Machine `json:"machine_information,omitempty"`
	GPUDriverVersion         *string          `json:"gpu_driver,omitempty"`
	GPUDevice                *int             `json:"gpu_device,omitempty"`
	GPUInformation           *nvidiasmi.GPU   `json:"gpu_information,omitempty"`
	InterconnectName         string           `json:"interconnect_name,omitempty"`
	TheoreticalGFlops        int64            `json:"theoretical_glops,omitempty"`
	MemoryBandwidth          float64          `json:"memory_bandwidth,omitempty"`
	IdealArithmeticIntensity float64          `json:"ideal_arithmetic_intensity,omitempty"`
	InterconnectBandwidth    float64          `json:"interconnect_bandwidth,omitempty"`
}

func (SummaryBase) Header(opts ...writer.Option) []string {
	return []string{
		"id",
		"created_at",
		"model_name",
		"framework_name",
		"using_gpu",
		"batch_size",
		"trace_level",
		"hostname",
		"machine_architecture",
		"machine_memory (KB)",
		"CPU_name",
		"GPU_name",
		"GPU_driver",
		"interconnect_name",
		"theoretical_flops (GFLOPS)",
		"memory_bandwidth (GB/s)",
		"interconnect_bandwidth (GB/s)",
		"ideal_arithmetic_intensity (flops/byte)",
	}
}

func (s SummaryBase) Row(opts ...writer.Option) []string {
	return []string{
		fmt.Sprintf(`%x`, string(s.ID)),
		s.CreatedAt.String(),
		s.ModelName,
		s.FrameworkName,
		cast.ToString(s.UsingGPU),
		cast.ToString(s.BatchSize),
		s.TraceLevel,
		s.HostName,
		s.MachineArchitecture,
		cast.ToString(s.MachineInformation.MemTotal),
		s.MachineInformation.CPU[0].ModelName,
		s.GPUInformation.ProductName,
		*s.GPUDriverVersion,
		s.InterconnectName,
		cast.ToString(s.TheoreticalGFlops),
		cast.ToString(s.MemoryBandwidth),
		cast.ToString(s.InterconnectBandwidth),
		fmt.Sprintf("%.2f", s.IdealArithmeticIntensity),
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
	theoreticalGFlops, err := e.GPUInformation.TheoreticalGFlops()
	if err != nil {
		log.WithError(err).Error("unable to get theoretical gflops")
	}
	memoryBandwidth, err := e.GPUInformation.MemoryBandwidth()
	if err != nil {
		log.WithError(err).Error("unable to get memory bandwidth")
	}
	interconnectName, err := e.GPUInformation.InterconnectName()
	if err != nil {
		log.WithError(err).Error("unable to get interconnect name")
	}
	interconnectBandwidth, err := e.GPUInformation.InterconnectBandwidth()
	if err != nil {
		log.WithError(err).Error("unable to get interconnect bandwidth")
	}

	idealArithmeticIntensity := float64(theoreticalGFlops) / memoryBandwidth

	return SummaryBase{
		ID:                       e.ID,
		CreatedAt:                e.CreatedAt,
		UpdatedAt:                time.Now(),
		ModelName:                e.Model.Name,
		ModelVersion:             e.Model.Version,
		FrameworkName:            e.Framework.Name,
		FrameworkVersion:         e.Framework.Version,
		MachineArchitecture:      e.MachineArchitecture,
		UsingGPU:                 e.UsingGPU,
		BatchSize:                e.BatchSize,
		HostName:                 e.Hostname,
		TraceLevel:               e.TraceLevel,
		MachineInformation:       e.MachineInformation,
		GPUDriverVersion:         e.GPUDriverVersion,
		GPUDevice:                e.GPUDevice,
		GPUInformation:           e.GPUInformation,
		TheoreticalGFlops:        theoreticalGFlops,
		MemoryBandwidth:          memoryBandwidth,
		IdealArithmeticIntensity: idealArithmeticIntensity,
		InterconnectBandwidth:    interconnectBandwidth,
		InterconnectName:         interconnectName,
	}
}

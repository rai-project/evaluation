package evaluation

import (
	"errors"

	"github.com/rai-project/evaluation/writer"
	"github.com/spf13/cast"
)

//easyjson:json
type SummaryGPUKernelModelAggreInformation struct {
	SummaryModelInformation `json:",inline"`
	Duration                float64 `json:"gpu_duration,omitempty"`
	Flops                   float64 `json:"flops,omitempty"`
	DramReadBytes           float64 `json:"dram_read_bytes,omitempty"`
	DramWriteBytes          float64 `json:"dram_write_bytes,omitempty"`
	ArithmeticIntensity     float64 `json:"arithmetic_intensity,omitempty"`
	ArithmeticThroughput    float64 `json:"arithmetic_throughput,omitempty"`
	MemoryBound             bool    `json:"memory_bound,omitempty"`
}

type SummaryGPUKernelModelAggreInformations []SummaryGPUKernelModelAggreInformation

func (p SummaryGPUKernelModelAggreInformations) Len() int { return len(p) }
func (p SummaryGPUKernelModelAggreInformations) Less(i, j int) bool {
	x := p[i]
	y := p[j]
	return x.Duration > y.Duration
}
func (p SummaryGPUKernelModelAggreInformations) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (info SummaryGPUKernelModelAggreInformation) Header(opts ...writer.Option) []string {
	return []string{
		"model_duration (us)",
		"model_gpu_duration (us)",
		"model_flops",
		"model_dram_read_bytes",
		"model_dram_write_bytes",
		"model_arithmetic_intensity (flops/byte)",
		"model_arithmetic_throughput (GFlops)",
		"model_memory_bound",
	}
}

func (info SummaryGPUKernelModelAggreInformation) Row(opts ...writer.Option) []string {
	return []string{
		cast.ToString(info.SummaryModelInformation.Duration),
		cast.ToString(info.Duration),
		cast.ToString(info.Flops),
		cast.ToString(info.DramReadBytes),
		cast.ToString(info.DramWriteBytes),
		cast.ToString(info.ArithmeticIntensity),
		cast.ToString(info.ArithmeticThroughput),
		cast.ToString(info.MemoryBound),
	}
}

func (es Evaluations) SummaryGPUKernelModelAggreInformations(perfCol *PerformanceCollection) (SummaryGPUKernelModelAggreInformations, error) {
	summary := SummaryGPUKernelModelAggreInformations{}
	gpuLayerInfos, err := es.SummaryGPUKernelLayerInformations(perfCol)
	if err != nil {
		return summary, errors.New("no span is found for the evaluation")
	}
	duration := float64(0)
	flops := float64(0)
	readBytes := float64(0)
	writeBytes := float64(0)
	for _, gpuLayerInfo := range gpuLayerInfos {
		if gpuLayerInfo.Index == 0 {
			continue
		}
		gpuInfos := gpuLayerInfo.SummaryGPUKernelInformations
		for _, gpuInfo := range gpuInfos {
			duration += gpuInfo.Duration
			flops += gpuInfo.MeanFlops
			readBytes += gpuInfo.MeanDramReadBytes
			writeBytes += gpuInfo.MeanDramWriteBytes
		}
	}

	modelInfos, err := (es.SummaryModelInformations(perfCol))
	modelInfo := modelInfos[0]
	if err != nil {
		modelInfo = SummaryModelInformation{}
	}

	arithmeticIntensity := flops / (readBytes + writeBytes)
	memoryBound := false
	if arithmeticIntensity < modelInfo.IdealArithmeticIntensity {
		memoryBound = true
	}
	arithmeticThroughput := flops / duration / float64(1000)

	summary = append(summary, SummaryGPUKernelModelAggreInformation{
		SummaryModelInformation: modelInfo,
		Duration:                duration,
		Flops:                   flops,
		DramReadBytes:           readBytes,
		DramWriteBytes:          writeBytes,
		ArithmeticIntensity:     arithmeticIntensity,
		ArithmeticThroughput:    arithmeticThroughput,
		MemoryBound:             memoryBound,
	})

	return summary, nil
}

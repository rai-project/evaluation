package evaluation

import (
	"errors"

	"github.com/rai-project/evaluation/writer"
	"github.com/spf13/cast"
)

//easyjson:json
type SummayGPULayerAggregatedInformation struct {
	SummaryLayerInformation `json:",inline"`
	GPUDuration             float64 `json:"gpu_duration,omitempty"`
	CPUDuration             float64 `json:"cpu_duration,omitempty"`
	Flops                   float64 `json:"flops,omitempty"`
	DramReadBytes           float64 `json:"dram_read_bytes,omitempty"`
	DramWriteBytes          float64 `json:"dram_write_bytes,omitempty"`
}

//easyjson:json
type SummayGPULayerAggregatedInformations []SummayGPULayerAggregatedInformation

func (SummayGPULayerAggregatedInformation) Header(iopts ...writer.Option) []string {
	extra := []string{
		"layer_index",
		"layer_name",
		"layer_type",
		"layer_duration (us)",
		"gpu_duration (us)",
		"cpu_duration (us)",
		"flops",
		"dram_read_bytes",
		"dram_write_bytes",
	}
	opts := writer.NewOptions(iopts...)
	if opts.ShowSummaryBase {
		return append(SummaryBase{}.Header(iopts...), extra...)
	}
	return extra
}

func (s SummayGPULayerAggregatedInformation) Row(iopts ...writer.Option) []string {
	extra := []string{
		cast.ToString(s.Index),
		s.Name,
		s.Type,
		cast.ToString(s.MeanDuration),
		cast.ToString(s.GPUDuration),
		cast.ToString(s.CPUDuration),
		cast.ToString(s.Flops),
		cast.ToString(s.DramReadBytes),
		cast.ToString(s.DramWriteBytes),
	}
	opts := writer.NewOptions(iopts...)
	if opts.ShowSummaryBase {
		return append(s.SummaryBase.Row(iopts...), extra...)
	}
	return extra
}

func (es Evaluations) SummayGPULayerAggregatedInformations(perfCol *PerformanceCollection) (SummayGPULayerAggregatedInformations, error) {
	summary := SummayGPULayerAggregatedInformations{}
	gpuLayerInfos, err := es.SummaryGPULayerInformations(perfCol)
	if err != nil {
		return summary, errors.New("no span is found for the evaluation")
	}

	for _, gpuLayerInfo := range gpuLayerInfos {
		layerInfo := gpuLayerInfo.SummaryLayerInformation
		gpuInfos := gpuLayerInfo.SummaryGPUInformations
		duration := float64(0)
		flops := float64(0)
		readBytes := float64(0)
		writeBytes := float64(0)
		for _, gpuInfo := range gpuInfos {
			duration += gpuInfo.MeanDuration
			flops += gpuInfo.MeanFlops
			readBytes += gpuInfo.MeanDramReadBytes
			writeBytes += gpuInfo.MeanDramWriteBytes
		}

		summary = append(summary, SummayGPULayerAggregatedInformation{
			SummaryLayerInformation: layerInfo,
			GPUDuration:             duration,
			CPUDuration:             layerInfo.MeanDuration - duration,
			Flops:                   flops,
			DramReadBytes:           readBytes,
			DramWriteBytes:          writeBytes,
		})
	}

	return summary, nil
}

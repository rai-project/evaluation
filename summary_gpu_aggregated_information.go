package evaluation

import (
	"github.com/rai-project/evaluation/writer"
	"github.com/spf13/cast"
)

type SummaryGPUAggregatedInformation struct {
	SummaryModelInformation `json:",inline"`
	Name                    string  `json:"name,omitempty"`
	Duration                float64 `json:"gpu_duration,omitempty"`
	Flops                   float64 `json:"flops,omitempty"`
	DramReadBytes           float64 `json:"dram_read_bytes,omitempty"`
	DramWriteBytes          float64 `json:"dram_write_bytes,omitempty"`
}

type SummaryGPUAggregatedInformations []SummaryGPUAggregatedInformation

func (p SummaryGPUAggregatedInformations) Len() int { return len(p) }
func (p SummaryGPUAggregatedInformations) Less(i, j int) bool {
	x := p[i]
	y := p[j]
	return x.Duration > y.Duration
}
func (p SummaryGPUAggregatedInformations) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (info SummaryGPUAggregatedInformation) Header(opts ...writer.Option) []string {
	return []string{
		"kernel_name",
		"kernel_duration (us)",
		"kernel_flops",
		"kernel_dram_read_bytes",
		"kernel_dram_write_bytes",
		"model_duration_percentage",
	}
}

func (info SummaryGPUAggregatedInformation) Row(opts ...writer.Option) []string {
	return []string{
		info.Name,
		cast.ToString(info.Duration),
		cast.ToString(info.Flops),
		cast.ToString(info.DramReadBytes),
		cast.ToString(info.DramWriteBytes),
		cast.ToString(info.Duration * float64(100) / info.SummaryModelInformation.Duration),
	}
}

func (es Evaluations) SummaryGPUAggregatedInformations(perfCol *PerformanceCollection) (SummaryGPUAggregatedInformations, error) {
	summary := SummaryGPUAggregatedInformations{}

	modelSummary, err := es.SummaryModelInformation(perfCol)
	if err != nil {
		modelSummary = SummaryModelInformation{}
	}

	cudaKernelInfos, err := es.GPUInformationSummary(perfCol)
	if err != nil {
		return summary, err
	}

	infoMap := make(map[string]SummaryGPUAggregatedInformation)
	for _, cki := range cudaKernelInfos {
		v, ok := infoMap[cki.Name]
		if !ok {
			infoMap[cki.Name] = SummaryGPUAggregatedInformation{
				SummaryModelInformation: modelSummary,
				Name:                    cki.Name,
				Duration:                cki.MeanDuration,
				Flops:                   cki.MeanFlops,
				DramReadBytes:           cki.MeanDramReadBytes,
				DramWriteBytes:          cki.MeanDramWriteBytes,
			}
		} else {
			v.Duration += cki.MeanDuration
			v.Flops += cki.MeanFlops
			v.DramReadBytes += cki.MeanDramReadBytes
			v.DramWriteBytes += cki.MeanDramWriteBytes
			v.SummaryModelInformation = modelSummary
			infoMap[cki.Name] = v
		}
	}

	for _, v := range infoMap {
		summary = append(summary, v)
	}

	return summary, nil
}

package evaluation

import (
	"errors"

	"github.com/rai-project/evaluation/writer"
	"github.com/spf13/cast"
)

type SummaryGPUKernelModelAggreInformation struct {
	SummaryModelInformation `json:",inline"`
	Name                    string  `json:"name,omitempty"`
	Duration                float64 `json:"gpu_duration,omitempty"`
	Flops                   float64 `json:"flops,omitempty"`
	DramReadBytes           float64 `json:"dram_read_bytes,omitempty"`
	DramWriteBytes          float64 `json:"dram_write_bytes,omitempty"`
	ArithmeticIntensity     float64 `json:"dram_write_bytes,omitempty"`
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
		"kernel_name",
		"kernel_duration (us)",
		"model_duration_percentage",
		"kernel_flops",
		"kernel_dram_read_bytes",
		"kernel_dram_write_bytes",
		"arithmetic_intensity (flops/byte)",
		"memory_bound",
	}
}

func (info SummaryGPUKernelModelAggreInformation) Row(opts ...writer.Option) []string {
	return []string{
		info.Name,
		cast.ToString(info.Duration),
		cast.ToString(info.Duration * float64(100) / info.SummaryModelInformation.Duration),
		cast.ToString(info.Flops),
		cast.ToString(info.DramReadBytes),
		cast.ToString(info.DramWriteBytes),
		cast.ToString(info.ArithmeticIntensity),
		cast.ToString(info.MemoryBound),
	}
}

func (es Evaluations) SummaryGPUKernelModelAggreInformations(perfCol *PerformanceCollection) (SummaryGPUKernelModelAggreInformations, error) {
	summary := SummaryGPUKernelModelAggreInformations{}

	if len(es.GroupByBatchSize()) != 1 {
		return summary, errors.New("evaluations do not exsit or are not with the same batch size")
	}

	modelInfos, err := (es.SummaryModelInformations(perfCol))
	modelInfo := modelInfos[0]
	if err != nil {
		modelInfo = SummaryModelInformation{}
	}

	infos := SummaryGPUKernelInformations{}
	gpuKernelLayerInfos, err := es.SummaryGPUKernelLayerInformations(perfCol)
	if err != nil {
		return summary, err
	}
	for _, v := range gpuKernelLayerInfos {
		infos = append(infos, v.SummaryGPUKernelInformations...)
	}

	infoMap := make(map[string]SummaryGPUKernelModelAggreInformation)
	for _, info := range infos {
		v, ok := infoMap[info.Name]
		if !ok {
			infoMap[info.Name] = SummaryGPUKernelModelAggreInformation{
				SummaryModelInformation: modelInfo,
				Name:                    info.Name,
				Duration:                info.Duration,
				Flops:                   info.MeanFlops,
				DramReadBytes:           info.MeanDramReadBytes,
				DramWriteBytes:          info.MeanDramWriteBytes,
			}
		} else {
			v.Duration += info.Duration
			v.Flops += info.MeanFlops
			v.DramReadBytes += info.MeanDramReadBytes
			v.DramWriteBytes += info.MeanDramWriteBytes
			v.SummaryModelInformation = modelInfo
			infoMap[info.Name] = v
		}
	}
	for _, v := range infoMap {
		arithmeticIntensity := v.Flops / (v.DramReadBytes + v.DramWriteBytes)
		v.ArithmeticIntensity = arithmeticIntensity
		memoryBound := false
		if arithmeticIntensity < v.IdealArithmeticIntensity {
			memoryBound = true
		}
		v.MemoryBound = memoryBound
		summary = append(summary, v)
	}

	return summary, nil
}

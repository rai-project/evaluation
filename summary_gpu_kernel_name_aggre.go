package evaluation

import (
	"fmt"

	"github.com/rai-project/evaluation/writer"
	"github.com/spf13/cast"
)

//easyjson:json
type SummaryGPUKernelNameAggreInformation struct {
	SummaryModelInformation `json:",inline"`
	Name                    string  `json:"name,omitempty"`
	Count                   int     `json:"count,omitempty"`
	Duration                float64 `json:"gpu_duration,omitempty"`
	Flops                   float64 `json:"flops,omitempty"`
	DramReadBytes           float64 `json:"dram_read_bytes,omitempty"`
	DramWriteBytes          float64 `json:"dram_write_bytes,omitempty"`
	AchievedOccupancy       float64 `json:"achieved_occupancy,omitempty"`
	ArithmeticIntensity     float64 `json:"arithmetic_intensity,omitempty"`
	ArithmeticThroughput    float64 `json:"arithmetic_throughput,omitempty"`
	MemoryBound             bool    `json:"memory_bound,omitempty"`
}

type SummaryGPUKernelNameAggreInformations []SummaryGPUKernelNameAggreInformation

func (p SummaryGPUKernelNameAggreInformations) Len() int { return len(p) }
func (p SummaryGPUKernelNameAggreInformations) Less(i, j int) bool {
	x := p[i]
	y := p[j]
	return x.Duration > y.Duration
}
func (p SummaryGPUKernelNameAggreInformations) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (info SummaryGPUKernelNameAggreInformation) Header(opts ...writer.Option) []string {
	return []string{
		"kernel_name",
		"kernel_count",
		"kernel_duration (us)",
		"model_duration_percentage",
		"kernel_flops",
		"kernel_dram_read_bytes",
		"kernel_dram_write_bytes",
		"kernel_achieved_occupancy (%)",
		"kernel_arithmetic_intensity (flops/byte)",
		"kernel_arithmetic_throughput (GFlops)",
		"kernel_memory_bound",
	}
}

func (info SummaryGPUKernelNameAggreInformation) Row(opts ...writer.Option) []string {
	return []string{
		info.Name,
		cast.ToString(info.Count),
		fmt.Sprintf("%.2f", info.Duration),
		fmt.Sprintf("%.2f", float64(info.Duration*100)/float64(info.SummaryModelInformation.Duration)),
		cast.ToString(info.Flops),
		fmt.Sprintf("%.2f", info.DramReadBytes),
		fmt.Sprintf("%.2f", info.DramWriteBytes),
		fmt.Sprintf("%.2f", info.AchievedOccupancy*100),
		fmt.Sprintf("%.2f", info.ArithmeticIntensity),
		fmt.Sprintf("%.2f", info.ArithmeticThroughput),
		cast.ToString(info.MemoryBound),
	}
}

func (es Evaluations) SummaryGPUKernelNameAggreInformations(perfCol *PerformanceCollection) (SummaryGPUKernelNameAggreInformations, error) {
	summary := SummaryGPUKernelNameAggreInformations{}
	infos := SummaryGPUKernelInformations{}
	gpuKernelLayerInfos, err := es.SummaryGPUKernelLayerInformations(perfCol)
	if err != nil {
		return summary, err
	}
	for _, v := range gpuKernelLayerInfos {
		infos = append(infos, v.SummaryGPUKernelInformations...)
	}

	modelInfos, err := (es.SummaryModelInformations(perfCol))
	modelInfo := modelInfos[0]
	if err != nil {
		modelInfo = SummaryModelInformation{}
	}

	infoMap := make(map[string]SummaryGPUKernelNameAggreInformation)
	for _, info := range infos {
		v, ok := infoMap[info.Name]
		if !ok {
			infoMap[info.Name] = SummaryGPUKernelNameAggreInformation{
				SummaryModelInformation: modelInfo,
				Name:                    info.Name,
				Duration:                info.MeanDuration,
				Count:                   0,
				Flops:                   info.MeanFlops,
				DramReadBytes:           info.MeanDramReadBytes,
				DramWriteBytes:          info.MeanDramWriteBytes,
				AchievedOccupancy:       info.MeanDuration * info.MeanAchievedOccupancy,
			}
		} else {
			v.Duration += info.MeanDuration
			v.Count += 1
			v.Flops += info.MeanFlops
			v.DramReadBytes += info.MeanDramReadBytes
			v.DramWriteBytes += info.MeanDramWriteBytes
			v.AchievedOccupancy += info.MeanDuration * info.MeanAchievedOccupancy
			v.SummaryModelInformation = modelInfo
			infoMap[info.Name] = v
		}
	}
	for _, v := range infoMap {
		v.ArithmeticIntensity = 0
		if (v.DramReadBytes + v.DramWriteBytes) != 0 {
			v.ArithmeticIntensity = v.Flops / (v.DramReadBytes + v.DramWriteBytes)
		}
		memoryBound := false
		if v.ArithmeticIntensity < v.IdealArithmeticIntensity {
			memoryBound = true
		}
		v.AchievedOccupancy = v.AchievedOccupancy / v.Duration
		v.MemoryBound = memoryBound
		v.ArithmeticThroughput = v.Flops / v.Duration / float64(1000)
		summary = append(summary, v)
	}

	return summary, nil
}

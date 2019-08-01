package evaluation

import (
	json "encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/rai-project/evaluation/writer"
	"github.com/rai-project/go-echarts/charts"
	"github.com/spf13/cast"
)

//easyjson:json
type SummaryGPUKernelLayerAggreInformation struct {
	SummaryLayerInformation `json:",inline"`
	GPUDuration             float64 `json:"gpu_duration,omitempty"`
	CPUDuration             float64 `json:"cpu_duration,omitempty"`
	Flops                   float64 `json:"flops,omitempty"`
	DramReadBytes           float64 `json:"dram_read_bytes,omitempty"`
	DramWriteBytes          float64 `json:"dram_write_bytes,omitempty"`
	AchievedOccupancy       float64 `json:"achieved_occupancy,omitempty"`
	ArithmeticIntensity     float64 `json:"arithmetic_intensity,omitempty"`
	ArithmeticThroughput    float64 `json:"arithmetic_throughput,omitempty"`
	MemoryBound             bool    `json:"memory_bound,omitempty"`
}

//easyjson:json
type SummaryGPUKernelLayerAggreInformations []SummaryGPUKernelLayerAggreInformation

func (p SummaryGPUKernelLayerAggreInformations) Len() int { return len(p) }
func (p SummaryGPUKernelLayerAggreInformations) Less(i, j int) bool {
	x := p[i]
	y := p[j]
	return x.Index < y.Index
}
func (p SummaryGPUKernelLayerAggreInformations) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (SummaryGPUKernelLayerAggreInformation) Header(iopts ...writer.Option) []string {
	extra := []string{
		"layer_index",
		"layer_name",
		"layer_type",
		"layer_duration (us)",
		"layer_gpu_duration (us)",
		"layer_cpu_duration (us)",
		"layer_flops",
		"layer_dram_read_bytes",
		"layer_dram_write_bytes",
		"layer_achieved_occupancy (%)",
		"layer_arithmetic_intensity (flops/byte)",
		"layer_arithmetic_throughput (GFlops)",
		"layer_memory_bound",
	}
	opts := writer.NewOptions(iopts...)
	if opts.ShowSummaryBase {
		return append(SummaryBase{}.Header(iopts...), extra...)
	}
	return extra
}

func (s SummaryGPUKernelLayerAggreInformation) Row(iopts ...writer.Option) []string {
	extra := []string{
		cast.ToString(s.Index),
		s.Name,
		s.Type,
		fmt.Sprintf("%.2f", s.Duration),
		fmt.Sprintf("%.2f", s.GPUDuration),
		fmt.Sprintf("%.2f", s.CPUDuration),
		cast.ToString(s.Flops),
		fmt.Sprintf("%.2f", s.DramReadBytes),
		fmt.Sprintf("%.2f", s.DramWriteBytes),
		fmt.Sprintf("%.2f", s.AchievedOccupancy*100),
		fmt.Sprintf("%.2f", s.ArithmeticIntensity),
		fmt.Sprintf("%.2f", s.ArithmeticThroughput),
		cast.ToString(s.MemoryBound),
	}
	opts := writer.NewOptions(iopts...)
	if opts.ShowSummaryBase {
		return append(s.SummaryBase.Row(iopts...), extra...)
	}
	return extra
}

func (es Evaluations) SummaryGPUKernelLayerAggreInformations(perfCol *PerformanceCollection) (SummaryGPUKernelLayerAggreInformations, error) {
	summary := SummaryGPUKernelLayerAggreInformations{}
	gpuLayerInfos, err := es.SummaryGPUKernelLayerInformations(perfCol)
	if err != nil {
		return summary, errors.New("no span is found for the evaluation")
	}

	for _, gpuLayerInfo := range gpuLayerInfos {
		if gpuLayerInfo.Index == 0 {
			continue
		}
		layerInfo := gpuLayerInfo.SummaryLayerInformation
		gpuInfos := gpuLayerInfo.SummaryGPUKernelInformations
		duration := float64(0)
		flops := float64(0)
		readBytes := float64(0)
		writeBytes := float64(0)
		achievedOccupancy := float64(0)
		for _, gpuInfo := range gpuInfos {
			duration += gpuInfo.MeanDuration
			flops += gpuInfo.MeanFlops
			readBytes += gpuInfo.MeanDramReadBytes
			writeBytes += gpuInfo.MeanDramWriteBytes
			achievedOccupancy += gpuInfo.MeanDuration * gpuInfo.MeanAchievedOccupancy
		}
		arithmeticIntensity := float64(0)
		if (readBytes + writeBytes) != 0 {
			arithmeticIntensity = flops / (readBytes + writeBytes)
		}
		memoryBound := false
		if arithmeticIntensity < layerInfo.IdealArithmeticIntensity {
			memoryBound = true
		}
		arithmeticThroughput := float64(0)
		if flops != 0 && duration != 0 {
			arithmeticThroughput = flops / duration / float64(1000)
		}
		if duration == 0 {
			achievedOccupancy = 0
		} else {
			achievedOccupancy = achievedOccupancy / duration
		}

		summary = append(summary, SummaryGPUKernelLayerAggreInformation{
			SummaryLayerInformation: layerInfo,
			GPUDuration:             duration,
			CPUDuration:             layerInfo.Duration - duration,
			Flops:                   flops,
			DramReadBytes:           readBytes,
			DramWriteBytes:          writeBytes,
			AchievedOccupancy:       achievedOccupancy,
			ArithmeticIntensity:     arithmeticIntensity,
			ArithmeticThroughput:    arithmeticThroughput,
			MemoryBound:             memoryBound,
		})
	}

	return summary, nil
}

type SummaryGPUKernelLayerFlopsInformations SummaryGPUKernelLayerAggreInformations

type SummaryGPUKernelLayerDramReadInformations SummaryGPUKernelLayerAggreInformations

type SummaryGPUKernelLayerDramWriteInformations SummaryGPUKernelLayerAggreInformations

type SummaryGPUKernelLayerLatencyInformations SummaryGPUKernelLayerAggreInformations

type SummaryGPUKernelLayerGPUCPUInformations SummaryGPUKernelLayerAggreInformations

type SummaryGPUKernelLayerAchievedOccupancyInformations SummaryGPUKernelLayerAggreInformations

func (o SummaryGPUKernelLayerFlopsInformations) PlotName() string {
	if len(o) == 0 {
		return ""
	}
	return o[0].ModelName + `
  Batch Size = ` + cast.ToString(o[0].BatchSize) + " GPU Kernel flops per Layer"
}

func (o SummaryGPUKernelLayerDramReadInformations) PlotName() string {
	if len(o) == 0 {
		return ""
	}
	return o[0].ModelName + `
  Batch Size = ` + cast.ToString(o[0].BatchSize) + " GPU Kernel Dram Read per Layer"
}

func (o SummaryGPUKernelLayerDramWriteInformations) PlotName() string {
	if len(o) == 0 {
		return ""
	}
	return o[0].ModelName + `
  Batch Size = ` + cast.ToString(o[0].BatchSize) + " GPU Kernel Dram Write per Layer"
}

func (o SummaryGPUKernelLayerLatencyInformations) PlotName() string {
	if len(o) == 0 {
		return ""
	}
	return o[0].ModelName + `
  Batch Size = ` + cast.ToString(o[0].BatchSize) + " GPU Kernel Latency per Layer"
}

func (o SummaryGPUKernelLayerGPUCPUInformations) PlotName() string {
	if len(o) == 0 {
		return ""
	}
	return o[0].ModelName + `
  Batch Size = ` + cast.ToString(o[0].BatchSize) + " GPU vs CPU Latency per Layer"
}

func (o SummaryGPUKernelLayerAchievedOccupancyInformations) PlotName() string {
	if len(o) == 0 {
		return ""
	}
	return o[0].ModelName + `
  Batch Size = ` + cast.ToString(o[0].BatchSize) + " GPU Kernel Achieved Occupany per Layer"
}

func (o SummaryGPUKernelLayerFlopsInformations) BarPlot() *charts.Bar {
	bar := charts.NewBar()
	bar = o.BarPlotAdd(bar)
	return bar
}

func (o SummaryGPUKernelLayerDramReadInformations) BarPlot() *charts.Bar {
	bar := charts.NewBar()
	bar = o.BarPlotAdd(bar)
	return bar
}

func (o SummaryGPUKernelLayerDramWriteInformations) BarPlot() *charts.Bar {
	bar := charts.NewBar()
	bar = o.BarPlotAdd(bar)
	return bar
}

func (o SummaryGPUKernelLayerLatencyInformations) BarPlot() *charts.Bar {
	bar := charts.NewBar()
	bar = o.BarPlotAdd(bar)
	return bar
}

func (o SummaryGPUKernelLayerGPUCPUInformations) BarPlot() *charts.Bar {
	bar := charts.NewBar()
	bar = o.BarPlotAdd(bar)
	return bar
}

func (o SummaryGPUKernelLayerAchievedOccupancyInformations) BarPlot() *charts.Bar {
	bar := charts.NewBar()
	bar = o.BarPlotAdd(bar)
	return bar
}

type GPUKernelLayerAggreInformationSelector func(elem SummaryGPUKernelLayerAggreInformation) float64

func (o SummaryGPUKernelLayerAggreInformations) barPlotAdd(bar *charts.Bar, elemSelector GPUKernelLayerAggreInformationSelector) *charts.Bar {
	labels := []string{}
	for _, elem := range o {
		labels = append(labels, elem.Name)
	}
	bar.AddXAxis(labels)

	data := make([]float64, len(o))
	for ii, elem := range o {
		data[ii] = elemSelector(elem)
	}
	bar.AddYAxis("", data)
	bar.SetSeriesOptions(
		charts.LabelTextOpts{Show: false},
		charts.TextStyleOpts{FontSize: DefaultSeriesFontSize},
	)

	jsLabelsBts, _ := json.Marshal(labels)
	jsFun := `function (name, index) {
	  var labels = ` + strings.Replace(string(jsLabelsBts), `"`, "'", -1) + `;
	  return labels.indexOf(name);
  }`
	bar.SetGlobalOptions(
		charts.XAxisOpts{Name: "Layer Index", Show: false, AxisLabel: charts.LabelTextOpts{Show: true, Formatter: charts.FuncOpts(jsFun)}},
	)

	return bar
}

func (o SummaryGPUKernelLayerFlopsInformations) BarPlotAdd(bar0 *charts.Bar) *charts.Bar {
	bar := SummaryGPUKernelLayerAggreInformations(o).barPlotAdd(bar0, func(elem SummaryGPUKernelLayerAggreInformation) float64 {
		return elem.Flops / float64(1000000000)
	})
	bar.SetGlobalOptions(
		charts.YAxisOpts{Name: "GFlops"},
	)
	return bar
}

func (o SummaryGPUKernelLayerDramReadInformations) BarPlotAdd(bar0 *charts.Bar) *charts.Bar {
	bar := SummaryGPUKernelLayerAggreInformations(o).barPlotAdd(bar0, func(elem SummaryGPUKernelLayerAggreInformation) float64 {
		return elem.DramReadBytes / float64(1024*1024)
	})
	bar.SetGlobalOptions(
		charts.YAxisOpts{Name: "MB"},
	)
	return bar
}

func (o SummaryGPUKernelLayerDramWriteInformations) BarPlotAdd(bar0 *charts.Bar) *charts.Bar {
	bar := SummaryGPUKernelLayerAggreInformations(o).barPlotAdd(bar0, func(elem SummaryGPUKernelLayerAggreInformation) float64 {
		return elem.DramWriteBytes / float64(1024*1024)
	})
	bar.SetGlobalOptions(
		charts.YAxisOpts{Name: "MB"},
	)
	return bar
}

func (o SummaryGPUKernelLayerLatencyInformations) BarPlotAdd(bar0 *charts.Bar) *charts.Bar {
	bar := SummaryGPUKernelLayerAggreInformations(o).barPlotAdd(bar0, func(elem SummaryGPUKernelLayerAggreInformation) float64 {
		return elem.Duration
	})
	bar.SetGlobalOptions(
		charts.YAxisOpts{Name: "us"},
	)
	return bar
}

func (o SummaryGPUKernelLayerGPUCPUInformations) BarPlotAdd(bar *charts.Bar) *charts.Bar {
	labels := []string{}
	for _, elem := range o {
		labels = append(labels, elem.Name)
	}
	bar.AddXAxis(labels)

	gpu := make([]float64, len(o))
	for ii, elem := range o {
		gpu[ii] = elem.GPUDuration
	}
	bar.AddYAxis("GPU", gpu, charts.BarOpts{Stack: "stack"})

	cpu := make([]float64, len(o))
	for ii, elem := range o {
		cpu[ii] = elem.CPUDuration
	}
	bar.AddYAxis("CPU", cpu, charts.BarOpts{Stack: "stack"})

	bar.SetSeriesOptions(
		charts.LabelTextOpts{Show: false},
		charts.TextStyleOpts{FontSize: DefaultSeriesFontSize},
	)

	jsLabelsBts, _ := json.Marshal(labels)
	jsFun := `function (name, index) {
	  var labels = ` + strings.Replace(string(jsLabelsBts), `"`, "'", -1) + `;
	  return labels.indexOf(name);
  }`
	bar.SetGlobalOptions(
		charts.LegendOpts{
			Left: "left",
		},
		charts.XAxisOpts{Name: "Layer Index", Show: false, AxisLabel: charts.LabelTextOpts{Show: true, Formatter: charts.FuncOpts(jsFun)}},
		charts.YAxisOpts{Name: "us"},
	)

	return bar
}

func (o SummaryGPUKernelLayerAchievedOccupancyInformations) BarPlotAdd(bar0 *charts.Bar) *charts.Bar {
	bar := SummaryGPUKernelLayerAggreInformations(o).barPlotAdd(bar0, func(elem SummaryGPUKernelLayerAggreInformation) float64 {
		return elem.AchievedOccupancy * float64(100)
	})
	bar.SetGlobalOptions(
		charts.YAxisOpts{Name: "%"},
	)
	return bar
}

func (o SummaryGPUKernelLayerFlopsInformations) WriteBarPlot(path string) error {
	return writeBarPlot(o, path)
}

func (o SummaryGPUKernelLayerDramReadInformations) WriteBarPlot(path string) error {
	return writeBarPlot(o, path)
}

func (o SummaryGPUKernelLayerDramWriteInformations) WriteBarPlot(path string) error {
	return writeBarPlot(o, path)
}

func (o SummaryGPUKernelLayerLatencyInformations) WriteBarPlot(path string) error {
	return writeBarPlot(o, path)
}

func (o SummaryGPUKernelLayerAchievedOccupancyInformations) WriteBarPlot(path string) error {
	return writeBarPlot(o, path)
}

func (o SummaryGPUKernelLayerGPUCPUInformations) WriteBarPlot(path string) error {
	return writeBarPlot(o, path)
}

func (o SummaryGPUKernelLayerFlopsInformations) OpenBarPlot() error {
	return openBarPlot(o)
}

func (o SummaryGPUKernelLayerDramReadInformations) OpenBarPlot() error {
	return openBarPlot(o)
}

func (o SummaryGPUKernelLayerDramWriteInformations) OpenBarPlot() error {
	return openBarPlot(o)
}

func (o SummaryGPUKernelLayerLatencyInformations) OpenBarPlot() error {
	return openBarPlot(o)
}

func (o SummaryGPUKernelLayerAchievedOccupancyInformations) OpenBarPlot() error {
	return openBarPlot(o)
}

func (o SummaryGPUKernelLayerGPUCPUInformations) OpenBarPlot() error {
	return openBarPlot(o)
}

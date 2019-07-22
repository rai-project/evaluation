package evaluation

import (
	"errors"
	"time"

	"github.com/rai-project/evaluation/writer"
	"github.com/rai-project/go-echarts/charts"
	"github.com/spf13/cast"
)

//easyjson:json
type SummayGPUKernelLayerAggreInformation struct {
	SummaryLayerInformation `json:",inline"`
	GPUDuration             float64 `json:"gpu_duration,omitempty"`
	CPUDuration             float64 `json:"cpu_duration,omitempty"`
	Flops                   float64 `json:"flops,omitempty"`
	DramReadBytes           float64 `json:"dram_read_bytes,omitempty"`
	DramWriteBytes          float64 `json:"dram_write_bytes,omitempty"`
	ArithmeticIntensity     float64 `json:"dram_write_bytes,omitempty"`
	MemoryBound             bool    `json:"memory_bound,omitempty"`
}

//easyjson:json
type SummayGPUKernelLayerAggreInformations []SummayGPUKernelLayerAggreInformation

func (p SummayGPUKernelLayerAggreInformations) Len() int { return len(p) }
func (p SummayGPUKernelLayerAggreInformations) Less(i, j int) bool {
	x := p[i]
	y := p[j]
	return x.Index < y.Index
}
func (p SummayGPUKernelLayerAggreInformations) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (SummayGPUKernelLayerAggreInformation) Header(iopts ...writer.Option) []string {
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
		"arithmetic_intensity (flops/byte)",
		"memory_bound",
	}
	opts := writer.NewOptions(iopts...)
	if opts.ShowSummaryBase {
		return append(SummaryBase{}.Header(iopts...), extra...)
	}
	return extra
}

func (s SummayGPUKernelLayerAggreInformation) Row(iopts ...writer.Option) []string {
	extra := []string{
		cast.ToString(s.Index),
		s.Name,
		s.Type,
		cast.ToString(s.Duration),
		cast.ToString(s.GPUDuration),
		cast.ToString(s.CPUDuration),
		cast.ToString(s.Flops),
		cast.ToString(s.DramReadBytes),
		cast.ToString(s.DramWriteBytes),
		cast.ToString(s.ArithmeticIntensity),
		cast.ToString(s.MemoryBound),
	}
	opts := writer.NewOptions(iopts...)
	if opts.ShowSummaryBase {
		return append(s.SummaryBase.Row(iopts...), extra...)
	}
	return extra
}

func (es Evaluations) SummaryGPUKernelLayerAggreInformations(perfCol *PerformanceCollection) (SummayGPUKernelLayerAggreInformations, error) {
	summary := SummayGPUKernelLayerAggreInformations{}
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
		for _, gpuInfo := range gpuInfos {
			duration += gpuInfo.Duration
			flops += gpuInfo.MeanFlops
			readBytes += gpuInfo.MeanDramReadBytes
			writeBytes += gpuInfo.MeanDramWriteBytes
		}
		arithmeticIntensity := flops / (readBytes + writeBytes)
		memoryBound := false
		if arithmeticIntensity < layerInfo.IdealArithmeticIntensity {
			memoryBound = true
		}
		summary = append(summary, SummayGPUKernelLayerAggreInformation{
			SummaryLayerInformation: layerInfo,
			GPUDuration:             duration,
			CPUDuration:             layerInfo.Duration - duration,
			Flops:                   flops,
			DramReadBytes:           readBytes,
			DramWriteBytes:          writeBytes,
			ArithmeticIntensity:     arithmeticIntensity,
			MemoryBound:             memoryBound,
		})
	}

	return summary, nil
}

type SummaryGPUKernelLayerGPUCPUInformations SummayGPUKernelLayerAggreInformations

type SummaryGPUKernelLayerFlopsInformations SummayGPUKernelLayerAggreInformations

type SummaryGPUKernelLayerDramReadInformations SummayGPUKernelLayerAggreInformations

type SummaryGPUKernelLayerDramWriteInformations SummayGPUKernelLayerAggreInformations

type SummaryGPUKernelLayerDurationInformations SummayGPUKernelLayerAggreInformations

func (o SummaryGPUKernelLayerFlopsInformations) PlotName() string {
	if len(o) == 0 {
		return ""
	}
	return o[0].ModelName + " Layer GPU Kernel Flops"
}

func (o SummaryGPUKernelLayerDramReadInformations) PlotName() string {
	if len(o) == 0 {
		return ""
	}
	return o[0].ModelName + " Layer GPU Kernel Dram Read Bytes"
}

func (o SummaryGPUKernelLayerDramWriteInformations) PlotName() string {
	if len(o) == 0 {
		return ""
	}
	return o[0].ModelName + " Layer GPU Kernel Dram Write Bytes"
}

func (o SummaryGPUKernelLayerDurationInformations) PlotName() string {
	if len(o) == 0 {
		return ""
	}
	return o[0].ModelName + " Layer GPU Kernel Duration"
}

func (o SummaryGPUKernelLayerFlopsInformations) BarPlot(title string) *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.TitleOpts{Title: title},
		charts.ToolboxOpts{Show: true},
	)
	bar = o.BarPlotAdd(bar)
	return bar
}

func (o SummaryGPUKernelLayerDramReadInformations) BarPlot(title string) *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.TitleOpts{Title: title},
		charts.ToolboxOpts{Show: true},
	)
	bar = o.BarPlotAdd(bar)
	return bar
}

func (o SummaryGPUKernelLayerDramWriteInformations) BarPlot(title string) *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.TitleOpts{Title: title},
		charts.ToolboxOpts{Show: true},
	)
	bar = o.BarPlotAdd(bar)
	return bar
}

func (o SummaryGPUKernelLayerDurationInformations) BarPlot(title string) *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.TitleOpts{Title: title},
		charts.ToolboxOpts{Show: true},
	)
	bar = o.BarPlotAdd(bar)
	return bar
}

type GPUKernelLayerAggreInformationSelector func(elem SummayGPUKernelLayerAggreInformation) float64

func (o SummayGPUKernelLayerAggreInformations) barPlotAdd(bar *charts.Bar, elemSelector GPUKernelLayerAggreInformationSelector) *charts.Bar {
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
	bar.SetSeriesOptions(charts.LabelTextOpts{Show: false})
	bar.SetGlobalOptions(
		charts.XAxisOpts{Name: "Layer Index"},
	)
	return bar
}

func (o SummaryGPUKernelLayerGPUCPUInformations) BarPlotAdd(bar0 *charts.Bar) *charts.Bar {
	bar := SummayGPUKernelLayerAggreInformations(o).barPlotAdd(bar0, func(elem SummayGPUKernelLayerAggreInformation) float64 {
		return elem.Flops
	})
	bar.SetGlobalOptions(
		charts.YAxisOpts{Name: "Latency(" + unitName(time.Microsecond) + ")"},
	)
	return bar
}

func (o SummaryGPUKernelLayerFlopsInformations) BarPlotAdd(bar0 *charts.Bar) *charts.Bar {
	bar := SummayGPUKernelLayerAggreInformations(o).barPlotAdd(bar0, func(elem SummayGPUKernelLayerAggreInformation) float64 {
		return elem.Flops / float64(1000000000)
	})
	bar.SetGlobalOptions(
		charts.YAxisOpts{Name: "GFlops"},
	)
	return bar
}

func (o SummaryGPUKernelLayerDramReadInformations) BarPlotAdd(bar0 *charts.Bar) *charts.Bar {
	bar := SummayGPUKernelLayerAggreInformations(o).barPlotAdd(bar0, func(elem SummayGPUKernelLayerAggreInformation) float64 {
		return elem.DramReadBytes / float64(1024*1024)
	})
	bar.SetGlobalOptions(
		charts.YAxisOpts{Name: "MB"},
	)
	return bar
}

func (o SummaryGPUKernelLayerDramWriteInformations) BarPlotAdd(bar0 *charts.Bar) *charts.Bar {
	bar := SummayGPUKernelLayerAggreInformations(o).barPlotAdd(bar0, func(elem SummayGPUKernelLayerAggreInformation) float64 {
		return elem.DramWriteBytes / float64(1024*1024)
	})
	bar.SetGlobalOptions(
		charts.YAxisOpts{Name: "MB"},
	)
	return bar
}

func (o SummaryGPUKernelLayerDurationInformations) BarPlotAdd(bar0 *charts.Bar) *charts.Bar {
	bar := SummayGPUKernelLayerAggreInformations(o).barPlotAdd(bar0, func(elem SummayGPUKernelLayerAggreInformation) float64 {
		return elem.Duration
	})
	bar.SetGlobalOptions(
		charts.YAxisOpts{Name: "us"},
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

func (o SummaryGPUKernelLayerDurationInformations) WriteBarPlot(path string) error {
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

func (o SummaryGPUKernelLayerDurationInformations) OpenBarPlot() error {
	return openBarPlot(o)
}

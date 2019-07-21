package evaluation

import (
	json "encoding/json"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/k0kubun/pp"
	"github.com/rai-project/evaluation/writer"
	"github.com/rai-project/go-echarts/charts"
	"github.com/rai-project/tracer"
	"github.com/spf13/cast"
	model "github.com/uber/jaeger/model/json"
)

//easyjson:json
type SummaryLayerInformation struct {
	SummaryModelInformation  `json:",inline"`
	Index                    int     `json:"index,omitempty"`
	Name                     string  `json:"layer_name,omitempty"`
	Type                     string  `json:"type,omitempty"`
	StaticType               string  `json:"static_type,omitempty"`
	Shape                    string  `json:"shap,omitempty"`
	Duration                 float64 `json:"mean_duration,omitempty"`
	Durations                []int64 `json:"durations,omitempty"`
	AllocatedBytes           []int64 `json:"allocated_bytes,omitempty"`      // Total number of bytes allocated if known
	PeakAllocatedBytes       []int64 `json:"peak_allocated_bytes,omitempty"` // Total number of ebytes allocated if known
	AllocatorBytesInUse      []int64 `json:"allocator_bytes_in_use,omitempty"`
	AllocatorName            string  `json:"allocator_name,omitempty"` // Name of the allocator used
	HostTempMemSizes         []int64 `json:"host_temp_mem_sizes,omitempty"`
	DeviceTempMemSizes       []int64 `json:"device_temp_mem_sizes,omitempty"`
	HostPersistentMemSizes   []int64 `json:"host_persistent_mem_sizes,omitempty"`
	DevicePersistentMemSizes []int64 `json:"device_persistent_mem_sizes,omitempty"`
}

//easyjson:json
type SummaryLayerInformations []SummaryLayerInformation

//easyjson:json
type SummaryMeanLayerInformation SummaryLayerInformation

//easyjson:json
type SummaryLayerMemoryInformations SummaryLayerInformations

//easyjson:json
type SummaryLayerLatencyInformations SummaryLayerInformations

func (SummaryLayerInformation) Header(iopts ...writer.Option) []string {
	extra := []string{
		"layer_index",
		"layer_name",
		"layer_type",
		"layer_shape",
		"layer_mean_duration (us)",
		"layer_durations (us)",
		"layer_allocated_bytes",
		"layer_peak_allocated_bytes",
		"layer_allocator_bytes_in_use",
		"layer_allocator_name",
		"layer_host_temp_mem_size",
		"layer_device_temp_mem_size",
		"layer_host_persistent_mem_size",
		"layer_device_persistent_mem_size",
	}
	opts := writer.NewOptions(iopts...)
	if opts.ShowSummaryBase {
		return append(SummaryBase{}.Header(iopts...), extra...)
	}
	return extra
}

func (s SummaryLayerInformation) Row(iopts ...writer.Option) []string {
	extra := []string{
		cast.ToString(s.Index),
		s.Name,
		s.Type,
		s.Shape,
		cast.ToString(s.Duration),
		strings.Join(int64SliceToStringSlice(s.Durations), DefaultDimiter),
		strings.Join(int64SliceToStringSlice(s.AllocatedBytes), DefaultDimiter),
		strings.Join(int64SliceToStringSlice(s.PeakAllocatedBytes), DefaultDimiter),
		strings.Join(int64SliceToStringSlice(s.AllocatorBytesInUse), DefaultDimiter),
		s.AllocatorName,
		strings.Join(int64SliceToStringSlice(s.HostTempMemSizes), DefaultDimiter),
		strings.Join(int64SliceToStringSlice(s.DeviceTempMemSizes), DefaultDimiter),
		strings.Join(int64SliceToStringSlice(s.HostPersistentMemSizes), DefaultDimiter),
		strings.Join(int64SliceToStringSlice(s.DevicePersistentMemSizes), DefaultDimiter),
	}
	opts := writer.NewOptions(iopts...)
	if opts.ShowSummaryBase {
		return append(s.SummaryBase.Row(iopts...), extra...)
	}
	return extra
}

func (s SummaryMeanLayerInformation) Header(opts ...writer.Option) []string {
	return SummaryLayerInformation(s).Header(opts...)
}

func (s SummaryMeanLayerInformation) Row(opts ...writer.Option) []string {
	return []string{
		cast.ToString(s.Index),
		s.Name,
		s.Type,
		s.Shape,
		cast.ToString(s.Duration),
		strings.Join(int64SliceToStringSlice(s.Durations), DefaultDimiter),
		cast.ToString(TrimmedMeanInt64Slice(s.AllocatedBytes, DefaultTrimmedMeanFraction)),
		cast.ToString(TrimmedMeanInt64Slice(s.PeakAllocatedBytes, DefaultTrimmedMeanFraction)),
		s.AllocatorName,
		cast.ToString(TrimmedMeanInt64Slice(s.HostTempMemSizes, DefaultTrimmedMeanFraction)),
		cast.ToString(TrimmedMeanInt64Slice(s.DeviceTempMemSizes, DefaultTrimmedMeanFraction)),
		cast.ToString(TrimmedMeanInt64Slice(s.HostPersistentMemSizes, DefaultTrimmedMeanFraction)),
		cast.ToString(TrimmedMeanInt64Slice(s.DevicePersistentMemSizes, DefaultTrimmedMeanFraction)),
	}
}

func getLayerInfoFromLayerSpan(span model.Span) SummaryLayerInformation {
	layerInfo := SummaryLayerInformation{}
	idx, err := getTagValueAsString(span, "layer_sequence_index")
	if err != nil || idx == "" {
		return layerInfo
	}
	shape, _ := getTagValueAsString(span, "shape")
	staticType, _ := getTagValueAsString(span, "static_type")
	allocationDesc := getAllocationDescription(span)
	allocatorName := allocationDesc.AllocatorName
	allocationBytes := allocationDesc.AllocatedBytes
	peakAllocationBytes := []int64{}
	allocatorBytesInUse := []int64{}
	memoryUsed, exist := getTensorFlowAllocatorMemoryUsed(span)
	if exist {
		m, err := cast.ToInt64E(memoryUsed.PeakBytes)
		if err == nil {
			peakAllocationBytes = []int64{m}
		}
		m, err = cast.ToInt64E(memoryUsed.AllocatorBytesInUse)
		if err == nil {
			allocatorBytesInUse = []int64{m}
		}
	}
	hostTempMemSizes := []int64{}
	hostTempMemSize, err := getTagValueAsString(span, "temp_memory_size")
	if err == nil && hostTempMemSize != "" {
		m, err := cast.ToInt64E(hostTempMemSize)
		if err == nil {
			hostTempMemSizes = []int64{m}
		}
	}
	deviceTempMemSizes := []int64{}
	deviceTempMemSize, err := getTagValueAsString(span, "device_temp_memory_size")
	if err == nil && deviceTempMemSize != "" {
		m, err := cast.ToInt64E(deviceTempMemSize)
		if err == nil {
			deviceTempMemSizes = []int64{m}
		}
	}
	hostPersistentMemSizes := []int64{}
	hostPersistentMemSize, err := getTagValueAsString(span, "persistent_memory_size")
	if err == nil && hostPersistentMemSize != "" {
		m, err := cast.ToInt64E(hostPersistentMemSize)
		if err == nil {
			hostPersistentMemSizes = []int64{m}
		}
	}
	devicePersistentMemSizes := []int64{}
	devicePersistentMemSize, err := getTagValueAsString(span, "device_persistent_memory_size")
	if err == nil && devicePersistentMemSize != "" {
		m, err := cast.ToInt64E(devicePersistentMemSize)
		if err == nil {
			devicePersistentMemSizes = []int64{m}
		}
	}
	layerInfo = SummaryLayerInformation{
		Index:      cast.ToInt(idx),
		Name:       span.OperationName,
		Type:       getOpName(span),
		StaticType: staticType,
		Shape:      shape,
		Duration:   -1,
		Durations: []int64{
			cast.ToInt64(span.Duration),
		},
		AllocatedBytes: []int64{
			cast.ToInt64(allocationBytes),
		},
		PeakAllocatedBytes:       peakAllocationBytes,
		AllocatorBytesInUse:      allocatorBytesInUse,
		AllocatorName:            allocatorName,
		HostTempMemSizes:         hostTempMemSizes,
		DeviceTempMemSizes:       deviceTempMemSizes,
		HostPersistentMemSizes:   hostPersistentMemSizes,
		DevicePersistentMemSizes: devicePersistentMemSizes,
	}
	return layerInfo
}

func (es Evaluations) SummaryLayerInformations(perfCol *PerformanceCollection) (SummaryLayerInformations, error) {
	summary := SummaryLayerInformations{}
	spans, err := es.GetSpansFromPerformanceCollection(perfCol)
	if err != nil {
		return summary, err
	}
	if len(spans) == 0 {
		return summary, errors.New("no span is found for the evaluation")
	}

	if len(es) == 0 {
		return summary, errors.New("no evaluation is found in the database")
	}

	cPredictSpans := spans.FilterByOperationNameAndEvalTraceLevel("c_predict", tracer.FRAMEWORK_TRACE.String())
	groupedLayerSpans, err := getGroupedLayerSpansFromSpans(cPredictSpans, spans)
	if err != nil {
		return summary, err
	}
	numGroups := len(groupedLayerSpans)
	if numGroups == 0 {
		return summary, errors.New("no group of spans is found")
	}

	modelInfos, err := (es.SummaryModelInformations(perfCol))
	modelInfo := modelInfos[0]
	if err != nil {
		modelInfo = SummaryModelInformation{}
	}

	groupedLayerInfos := make([][]SummaryLayerInformation, numGroups)
	for ii, spans := range groupedLayerSpans {
		if groupedLayerInfos[ii] == nil {
			groupedLayerInfos[ii] = []SummaryLayerInformation{}
		}
		for _, span := range spans {
			if strings.HasPrefix(span.OperationName, "_") {
				continue
			}
			layerInfo := getLayerInfoFromLayerSpan(span)
			groupedLayerInfos[ii] = append(groupedLayerInfos[ii], layerInfo)
		}
	}

	for _, li := range groupedLayerInfos[0] {
		layerInfo := li
		for _, lis := range groupedLayerInfos[1:] {
			for _, lli := range lis {
				if lli.Name != li.Name || li.Index != li.Index {
					continue
				}
				layerInfo.Durations = append(layerInfo.Durations, lli.Durations...)
				layerInfo.AllocatedBytes = append(layerInfo.AllocatedBytes, lli.AllocatedBytes...)
				layerInfo.PeakAllocatedBytes = append(layerInfo.PeakAllocatedBytes, lli.PeakAllocatedBytes...)
				layerInfo.HostTempMemSizes = append(layerInfo.HostTempMemSizes, lli.HostTempMemSizes...)
				layerInfo.DeviceTempMemSizes = append(layerInfo.DeviceTempMemSizes, lli.DeviceTempMemSizes...)
				layerInfo.HostPersistentMemSizes = append(layerInfo.HostPersistentMemSizes, lli.HostPersistentMemSizes...)
				layerInfo.DevicePersistentMemSizes = append(layerInfo.DevicePersistentMemSizes, lli.DevicePersistentMemSizes...)
			}
		}

		layerInfo.SummaryModelInformation = modelInfo
		layerInfo.Duration = TrimmedMeanInt64Slice(layerInfo.Durations, DefaultTrimmedMeanFraction)
		summary = append(summary, layerInfo)
	}

	return summary, nil
}

func sortByLayerIndex(spans Spans) {
	sort.Slice(spans, func(ii, jj int) bool {
		li, foundI := spanTagValue(spans[ii], "layer_sequence_index")
		if !foundI {
			return false
		}
		lj, foundJ := spanTagValue(spans[jj], "layer_sequence_index")
		if !foundJ {
			return true
		}

		return cast.ToInt64(li) < cast.ToInt64(lj)
	})
}

func getGroupedLayerSpansFromSpans(cPredictSpans Spans, spans Spans) ([]Spans, error) {
	groupedSpans, err := getGroupedSpansFromSpans(cPredictSpans, spans)
	if err != nil {
		return nil, err
	}
	numPredictSpans := len(groupedSpans)

	groupedLayerSpans := make([]Spans, numPredictSpans)
	for ii, grsp := range groupedSpans {
		if len(grsp) == 0 {
			continue
		}

		groupedLayerSpans[ii] = Spans{}
		for _, sp := range grsp {
			traceLevel, err := getTagValueAsString(sp, "trace_level")
			if err != nil || traceLevel == "" {
				continue
			}
			if tracer.LevelFromName(traceLevel) != tracer.FRAMEWORK_TRACE {
				continue
			}
			groupedLayerSpans[ii] = append(groupedLayerSpans[ii], sp)
		}

		sortByLayerIndex(groupedLayerSpans[ii])
	}

	return groupedLayerSpans, nil
}

func (s SummaryLayerInformations) GetLayerInfoByName(name string) SummaryLayerInformation {
	for _, info := range s {
		if info.Name == name {
			return info
		}
	}
	return SummaryLayerInformation{}
}

func (o SummaryLayerLatencyInformations) PlotName() string {
	if len(o) == 0 {
		return ""
	}
	return o[0].ModelName + " Layer Latency"
}

func (o SummaryLayerMemoryInformations) PlotName() string {
	if len(o) == 0 {
		return ""
	}
	return o[0].ModelName + " Layer Allocated Memory"
}

func (o SummaryLayerLatencyInformations) BarPlot(title string) *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.TitleOpts{Title: title},
		charts.ToolboxOpts{Show: true},
	)
	bar = o.BarPlotAdd(bar)
	return bar
}

func (o SummaryLayerMemoryInformations) BarPlot(title string) *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.TitleOpts{Title: title},
		charts.ToolboxOpts{Show: true},
	)
	bar = o.BarPlotAdd(bar)
	return bar
}

type LayerInformationSelector func(elem SummaryLayerInformation) float64

func (o SummaryLayerInformations) barPlotAdd(bar *charts.Bar, elemSelector LayerInformationSelector) *charts.Bar {
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

func (o SummaryLayerLatencyInformations) BarPlotAdd(bar0 *charts.Bar) *charts.Bar {
	bar := SummaryLayerInformations(o).barPlotAdd(bar0, func(elem SummaryLayerInformation) float64 {
		return TrimmedMeanInt64Slice(elem.Durations, 0)
	})
	bar.SetGlobalOptions(
		charts.YAxisOpts{Name: "Latency(" + unitName(time.Microsecond) + ")"},
	)
	return bar
}

func (o SummaryLayerMemoryInformations) BarPlotAdd(bar0 *charts.Bar) *charts.Bar {
	bar := SummaryLayerInformations(o).barPlotAdd(bar0, func(elem SummaryLayerInformation) float64 {
		return TrimmedMeanInt64Slice(elem.AllocatedBytes, 0) / float64(1048576)
	})
	bar.SetGlobalOptions(
		charts.YAxisOpts{Name: "Allocated Memory(MB)"},
	)
	return bar

}

func (o SummaryLayerLatencyInformations) WriteBarPlot(path string) error {
	return writeBarPlot(o, path)
}

func (o SummaryLayerMemoryInformations) WriteBarPlot(path string) error {
	return writeBarPlot(o, path)
}

func (o SummaryLayerLatencyInformations) OpenBarPlot() error {
	return openBarPlot(o)
}

func (o SummaryLayerMemoryInformations) OpenBarPlot() error {
	return openBarPlot(o)
}

func (o SummaryLayerLatencyInformations) BoxPlot(title string) *charts.BoxPlot {
	box := charts.NewBoxPlot()
	box.SetGlobalOptions(
		charts.TitleOpts{Title: title},
		charts.ToolboxOpts{Show: true},
	)
	box = o.BoxPlotAdd(box)
	return box
}

func (o SummaryLayerLatencyInformations) BoxPlotAdd(box *charts.BoxPlot) *charts.BoxPlot {
	timeUnit := time.Microsecond

	isPrivate := func(info SummaryLayerInformation) bool {
		return strings.HasPrefix(info.Name, "_")
	}

	labels := []string{}
	for _, elem := range o {
		if isPrivate(elem) {
			continue
		}
		labels = append(labels, elem.Name)
	}
	box.AddXAxis(labels)

	durations := make([][]time.Duration, 0, len(o))
	for _, elem := range o {
		if isPrivate(elem) {
			continue
		}
		ts := make([]time.Duration, len(elem.Durations))
		for jj, t := range elem.Durations {
			ts[jj] = time.Duration(t)
		}
		durations = append(durations, prepareBoxplotData(ts))
	}
	if false {
		pp.Println(len(labels))
		pp.Println(len(durations))
		pp.Println(len(durations[0]))
	}
	box.AddYAxis("", durations)
	box.SetSeriesOptions(charts.LabelTextOpts{Show: false})
	jsLabelsBts, _ := json.Marshal(labels)
	jsFun := `function (name, index) {
    var labels = ` + strings.Replace(string(jsLabelsBts), `"`, "'", -1) + `;
    return labels.indexOf(name);
  }`
	box.SetGlobalOptions(
		charts.XAxisOpts{
			Name:      "Layer Name",
			Type:      "category",
			AxisLabel: charts.LabelTextOpts{Show: true, Rotate: 45, Formatter: charts.FuncOpts(jsFun)},
			SplitLine: charts.SplitLineOpts{Show: false},
			SplitArea: charts.SplitAreaOpts{Show: true},
		},
		charts.YAxisOpts{
			Name: "Latency(" + unitName(timeUnit) + ")",
			Type: "value",
			// NameRotate: 90,
			AxisLabel: charts.LabelTextOpts{Formatter: "{value}" + unitName(timeUnit)},
			SplitArea: charts.SplitAreaOpts{Show: true},
			Mix:       0,
		},
		charts.DataZoomOpts{
			Type:       "slider",
			XAxisIndex: []int{0},
			Start:      0,
			End:        float32(len(labels)),
		},
	)
	return box
}

func prepareBoxplotData(ds []time.Duration) []time.Duration {
	min := durationMin(ds)
	q1 := durationPercentile(ds, 25)
	q2 := durationPercentile(ds, 50)
	q3 := durationPercentile(ds, 75)
	max := durationMax(ds)
	return []time.Duration{min, q1, q2, q3, max}
}

func (o SummaryLayerLatencyInformations) WriteBoxPlot(path string) error {
	return writeBoxPlot(o, path)
}

func (o SummaryLayerLatencyInformations) OpenBoxPlot() error {
	return openBoxPlot(o)
}

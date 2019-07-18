package evaluation

import (
	"encoding/json"
	"errors"
	"sort"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/rai-project/evaluation/writer"
	"github.com/rai-project/tracer"
	trace_tree "github.com/rai-project/tracer/convert"
	"github.com/spf13/cast"
	model "github.com/uber/jaeger/model/json"
)

type Metadata map[string]interface{}

type SummaryGPUKernelInformation struct {
	Name               string     `json:"name,omitempty"`
	MangledName        string     `json:"mangled_name,omitempty"`
	Durations          []int64    `json:"durations,omitempty"`
	Tags               []Metadata `json:"tags,omitempty"`
	Logs               []Metadata `json:"logs,omitempty"`
	CorrelationId      int64      `json:"correlation_id,omitempty"`
	Duration           float64    `json:"mean_duration,omitempty"`
	MeanFlops          float64    `json:"mean_flops,omitempty"`
	MeanDramReadBytes  float64    `json:"mean_dram_read_bytes,omitempty"`
	MeanDramWriteBytes float64    `json:"mean_dram_write_bytes,omitempty"`
}

type SummaryGPUKernelInformations []SummaryGPUKernelInformation

func (p SummaryGPUKernelInformations) Len() int { return len(p) }
func (p SummaryGPUKernelInformations) Less(i, j int) bool {
	x := p[i]
	y := p[j]
	xDuration := TrimmedMeanInt64Slice(x.Durations, DefaultTrimmedMeanFraction)
	yDuration := TrimmedMeanInt64Slice(y.Durations, DefaultTrimmedMeanFraction)
	return xDuration > yDuration
}
func (p SummaryGPUKernelInformations) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (info SummaryGPUKernelInformation) Header(opts ...writer.Option) []string {
	extraHeader := []string{
		"kernel_name",
		"kernel_mean_duration (us)",
		"kernel_mean_flops",
		"kernel_mean_dram_read_bytes",
		"kernel_mean_dram_write_bytes",
		"kernel_durations (us)",
	}
	kernelLogKeys := SummaryGPUKernelInformations{info}.GetKernelLogKeys()
	if len(kernelLogKeys) != 0 {
		extraHeader = append(extraHeader, kernelLogKeys...)
	}
	return extraHeader
}

func (info SummaryGPUKernelInformation) Row(opts ...writer.Option) []string {
	extra := []string{
		info.Name,
		cast.ToString(info.Duration),
		cast.ToString(info.MeanFlops),
		cast.ToString(info.MeanDramReadBytes),
		cast.ToString(info.MeanDramWriteBytes),
		strings.Join(int64SliceToStringSlice(info.Durations), DefaultDimiter),
	}
	kernelLogKeys := SummaryGPUKernelInformations{info}.GetKernelLogKeys()
	for _, kernelLogKey := range kernelLogKeys {
		kernelLogs := []string{}
		for _, kernelLog := range info.Logs {
			for kernelLogKeyName, keryeLogValue := range kernelLog {
				if kernelLogKeyName == kernelLogKey {
					kernelLogs = append(kernelLogs, cast.ToString(keryeLogValue))
				}
			}
		}
		extra = append(extra, strings.Join(kernelLogs, DefaultDimiter))
	}
	return extra
}

type SummaryGPUKernelLayerInformation struct {
	SummaryLayerInformation      `json:",inline"`
	SummaryGPUKernelInformations SummaryGPUKernelInformations `json:"kernel_launch_information,omitempty"`
}

func (p SummaryGPUKernelLayerInformation) Len() int { return len(p.SummaryGPUKernelInformations) }
func (p SummaryGPUKernelLayerInformation) Less(i, j int) bool {
	x := p.SummaryGPUKernelInformations[i]
	y := p.SummaryGPUKernelInformations[j]
	xDuration := TrimmedMeanInt64Slice(x.Durations, DefaultTrimmedMeanFraction)
	yDuration := TrimmedMeanInt64Slice(y.Durations, DefaultTrimmedMeanFraction)
	return xDuration > yDuration
}
func (p SummaryGPUKernelLayerInformation) Swap(i, j int) {
	p.SummaryGPUKernelInformations[i], p.SummaryGPUKernelInformations[j] = p.SummaryGPUKernelInformations[j], p.SummaryGPUKernelInformations[i]
}

type SummaryGPUKernelLayerInformations []SummaryGPUKernelLayerInformation

func (p SummaryGPUKernelLayerInformations) Len() int { return len(p) }
func (p SummaryGPUKernelLayerInformations) Less(i, j int) bool {
	return p[i].SummaryLayerInformation.Index < p[j].SummaryLayerInformation.Index
}
func (p SummaryGPUKernelLayerInformations) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (infos SummaryGPUKernelLayerInformations) Header(opts ...writer.Option) []string {
	extraHeader := []string{
		"kernel_name",
		"kernel_durations (us)",
	}
	kernelLogKeys := infos.GetKernelLogKeys()
	if len(kernelLogKeys) != 0 {
		extraHeader = append(extraHeader, kernelLogKeys...)
	}
	return append(SummaryLayerInformation{}.Header(opts...), extraHeader...)
}

func (infos SummaryGPUKernelLayerInformations) Row(opts ...writer.Option) []string {
	return []string{}
}

func getMetaDataKeys(metadatas []Metadata) []string {
	if metadatas == nil || len(metadatas) == 0 {
		return []string{}
	}
	keyVisited := map[string]bool{}
	keys := []string{}
	for _, metadata := range metadatas {
		for key, _ := range metadata {
			if _, ok := keyVisited[key]; ok {
				continue
			}
			keys = append(keys, key)
			keyVisited[key] = true
		}
	}
	return keys
}

func (infos SummaryGPUKernelInformations) GetKernelLogKeys() []string {
	kernelLogs := []Metadata{}
	for _, info := range infos {
		if len(info.Logs) == 0 {
			continue
		}
		kernelLogs = append(kernelLogs, info.Logs...)
	}
	return getMetaDataKeys(kernelLogs)
}

func (infos SummaryGPUKernelLayerInformations) GetKernelLogKeys() []string {
	kernelLogs := []Metadata{}
	for _, info := range infos {
		for _, gpuKernelInformation := range info.SummaryGPUKernelInformations {
			if len(gpuKernelInformation.Logs) == 0 {
				continue
			}
			kernelLogs = append(kernelLogs, gpuKernelInformation.Logs...)
		}
	}
	return getMetaDataKeys(kernelLogs)
}

// Rows ...
func (info SummaryGPUKernelLayerInformation) Rows(iopts ...writer.Option) [][]string {
	gpuKernelInfos := info.SummaryGPUKernelInformations
	layerInfo := SummaryMeanLayerInformation(info.SummaryLayerInformation)
	layerInfoRow := layerInfo.Row(iopts...)

	opts := writer.NewOptions(iopts...)

	rows := [][]string{}

	kernelLogKeys := SummaryGPUKernelLayerInformations{info}.GetKernelLogKeys()

	isFilteredKernel := func(kernelInfo SummaryGPUKernelInformation) bool {
		if len(opts.FilterKernelNames) == 0 {
			return true
		}
		name := strings.ToLower(kernelInfo.MangledName)
		for _, filterName := range opts.FilterKernelNames {
			if name == strings.ToLower(filterName) {
				return true
			}
		}
		return false
	}

	for _, cki := range gpuKernelInfos {
		if !isFilteredKernel(cki) {
			continue
		}
		kernelTags, err := json.Marshal(cki.Tags)
		if err != nil {
			kernelTags = []byte{}
		}

		_ = kernelTags

		extra := []string{
			cki.Name,
			strings.Join(int64SliceToStringSlice(cki.Durations), DefaultDimiter),
		}

		for _, kernelLogKey := range kernelLogKeys {
			kernelLogs := []string{}
			for _, kernelLog := range cki.Logs {
				for kernelLogKeyName, keryeLogValue := range kernelLog {
					if kernelLogKeyName == kernelLogKey {
						kernelLogs = append(kernelLogs, cast.ToString(keryeLogValue))
					}
				}
			}
			extra = append(extra, strings.Join(kernelLogs, DefaultDimiter))
		}
		rows = append(rows, append(layerInfoRow, extra...))
	}
	return rows
}

func (k *SummaryGPUKernelInformation) addLogs(spanLogs []model.Log) {
	if k.Logs == nil {
		k.Logs = []Metadata{}
	}
	logs := Metadata{}
	for _, v := range spanLogs {
		for _, f := range v.Fields {
			logs[f.Key] = f.Value
		}
	}
	if len(logs) == 0 {
		return
	}
	k.Logs = append(k.Logs, logs)
}

func (k *SummaryGPUKernelInformation) addTags(spanTags []model.KeyValue) {
	if k.Tags == nil {
		k.Tags = []Metadata{}
	}
	tags := Metadata{}
	for _, v := range spanTags {
		tags[v.Key] = v.Value
	}
	if len(tags) == 0 {
		return
	}
	k.Tags = append(k.Tags, tags)
}

func GetMeanLogValue(info SummaryGPUKernelInformation, name string, trimmedMeanFraction float64) float64 {
	if info.Logs == nil {
		info.Logs = []Metadata{}
	}
	kernelLogs := []int64{}
	for _, kernelLog := range info.Logs {
		for kernelLogKeyName, keryeLogValue := range kernelLog {
			if kernelLogKeyName == name {
				kernelLogs = append(kernelLogs, cast.ToInt64(keryeLogValue))
			}
		}
	}
	return TrimmedMeanInt64Slice(kernelLogs, trimmedMeanFraction)
}

func GPUKernelSpantoGPUInformation(span model.Span) SummaryGPUKernelInformation {
	info := &SummaryGPUKernelInformation{
		Name:          mustGetTagValueAsString(span, "kernel_name"),
		MangledName:   mustGetTagValueAsString(span, "name"),
		Tags:          []Metadata{},
		Logs:          []Metadata{},
		CorrelationId: mustGetTagValueAsInt64(span, "correlation_id"),
		Durations: []int64{
			cast.ToInt64(span.Duration),
		},
	}
	info.addTags(span.Tags)
	info.addLogs(span.Logs)
	return *info
}

func CUDALaunchSpantoGPUInformation(span model.Span) SummaryGPUKernelInformation {
	kernelName := mustGetTagValueAsString(span, "kernel")
	info := &SummaryGPUKernelInformation{
		Name:          demangleName(kernelName),
		MangledName:   kernelName,
		Tags:          []Metadata{},
		Logs:          []Metadata{},
		CorrelationId: mustGetTagValueAsInt64(span, "correlation_id"),
	}
	info.addTags(span.Tags)
	info.addLogs(span.Logs)
	return *info
}

func (es Evaluations) SummaryGPUKernelLayerInformations(perfCol *PerformanceCollection) (SummaryGPUKernelLayerInformations, error) {
	summary := SummaryGPUKernelLayerInformations{}
	if len(es) == 0 {
		return summary, errors.New("no evaluation is found in the database")
	}

	layerInfos, err := es.SummaryLayerInformations(perfCol)
	if err != nil {
		layerInfos = SummaryLayerInformations{}
	}

	spans, err := es.GetSpansFromPerformanceCollection(perfCol)
	if err != nil {
		return summary, err
	}
	if len(spans) == 0 {
		return summary, errors.New("no span is found for the evaluation")
	}

	cPredictSpans := spans.FilterByOperationNameAndEvalTraceLevel("c_predict", tracer.SYSTEM_LIBRARY_TRACE.String())
	groupedSpans, err := getGroupedSpansFromSpans(cPredictSpans, spans)
	if err != nil {
		return summary, err
	}
	numGroups := len(groupedSpans)
	if numGroups == 0 {
		return summary, errors.New("no group of spans is found")
	}

	groupedLayerGPUInfos := make([][]SummaryGPUKernelLayerInformation, numGroups)
	for ii, grsp := range groupedSpans {
		if groupedLayerGPUInfos[ii] == nil {
			groupedLayerGPUInfos[ii] = []SummaryGPUKernelLayerInformation{}
		}

		trace := model.Trace{
			TraceID: "0",
			Spans:   grsp,
		}
		tree, err := trace_tree.NewIntervalTree(trace)
		if err != nil {
			panic(err)
		}

		for _, sp := range grsp {
			traceLevel, err := getTagValueAsString(sp, "trace_level")
			if err != nil || traceLevel == "" {
				continue
			}

			if tracer.LevelFromName(traceLevel) != tracer.FRAMEWORK_TRACE {
				continue
			}

			layerInterval := trace_tree.ToInterval(sp)
			layerSpan := *layerInterval.Span
			layerChildren := tree.ChildrenOf(layerInterval)

			layerInfo := SummaryLayerInformation{}
			if len(layerInfos) == 0 {
				idx, err := getTagValueAsString(layerSpan, "layer_sequence_index")
				if err != nil || idx == "" {
					return summary, errors.New("cannot find tag layer_sequence_index")
				}
				allocationDesc := getAllocationDescription(layerSpan)
				memoryUsed := getTensorFlowAllocatorMemoryUsed(layerSpan)
				allocationBytes := allocationDesc.AllocatedBytes
				peakAllocationBytes := memoryUsed.PeakBytes
				hostTempMemSize, _ := getTagValueAsString(layerSpan, "temp_memory_size")
				deviceTempMemSize, _ := getTagValueAsString(layerSpan, "device_temp_memory_size")
				hostPersistentMemSize, _ := getTagValueAsString(layerSpan, "persistent_memory_size")
				devicePersistentMemSize, _ := getTagValueAsString(layerSpan, "device_persistent_memory_size")
				layerInfo = SummaryLayerInformation{
					Index:     cast.ToInt(idx),
					Name:      layerSpan.OperationName,
					Type:      getOpName(layerSpan),
					Durations: []int64{},
					AllocatedBytes: []int64{
						cast.ToInt64(allocationBytes),
					},
					PeakAllocatedBytes: []int64{
						cast.ToInt64(peakAllocationBytes),
					},
					HostTempMemSizes: []int64{
						cast.ToInt64(hostTempMemSize),
					},
					DeviceTempMemSizes: []int64{
						cast.ToInt64(deviceTempMemSize),
					},
					HostPersistentMemSizes: []int64{
						cast.ToInt64(hostPersistentMemSize),
					},
					DevicePersistentMemSizes: []int64{
						cast.ToInt64(devicePersistentMemSize),
					},
				}
			} else {
				layerInfo = layerInfos.GetLayerInfoByName(layerSpan.OperationName)
			}

			layerGPUInformation := SummaryGPUKernelLayerInformation{
				SummaryLayerInformation:      layerInfo,
				SummaryGPUKernelInformations: []SummaryGPUKernelInformation{},
			}

			for _, childInterval := range layerChildren {
				child := *childInterval.Span
				traceLevel, err := getTagValueAsString(child, "trace_level")
				if err != nil || traceLevel == "" {
					continue
				}
				if tracer.LevelFromName(traceLevel) != tracer.SYSTEM_LIBRARY_TRACE {
					continue
				}
				if strings.ToLower(child.OperationName) != "cuda_launch" {
					continue
				}
				layerGPUInformation.SummaryGPUKernelInformations = append(layerGPUInformation.SummaryGPUKernelInformations, CUDALaunchSpantoGPUInformation(child))
			}

			for _, childInterval := range layerChildren {
				child := *childInterval.Span
				traceLevel, err := getTagValueAsString(child, "trace_level")
				if err != nil || traceLevel == "" {
					continue
				}
				if tracer.LevelFromName(traceLevel) != tracer.SYSTEM_LIBRARY_TRACE {
					continue
				}

				if strings.ToLower(child.OperationName) != "gpu_kernel" {
					continue
				}

				childCorrelationId, err := getTagValueAsInt64(child, "correlation_id")
				if err != nil {
					log.WithError(err).Error("expecting cuda launch to have a correlation_id")
					continue
				}
				for infoIdx := range layerGPUInformation.SummaryGPUKernelInformations {
					info := layerGPUInformation.SummaryGPUKernelInformations[infoIdx]
					if info.CorrelationId != childCorrelationId {
						continue
					}
					// only record kernel duration when no gpu metrics are captured
					if len(info.Logs) == 0 {
						info.Durations = []int64{
							cast.ToInt64(child.Duration),
						}
					}
					layerGPUInformation.SummaryGPUKernelInformations[infoIdx] = info
				}
			}
			groupedLayerGPUInfos[ii] = append(groupedLayerGPUInfos[ii], layerGPUInformation)
		}
	}

	layerGPUInfos := []SummaryGPUKernelLayerInformation{}
	for _, li := range groupedLayerGPUInfos[0] {
		layerGPUInfo := li
		for ii := range layerGPUInfo.SummaryGPUKernelInformations {
			cki := layerGPUInfo.SummaryGPUKernelInformations[ii]
			for _, lis := range groupedLayerGPUInfos[1:] {
				for _, lli := range lis {
					if lli.Name != li.Name || li.Index != li.Index {
						continue
					}
					for _, ccki := range lli.SummaryGPUKernelInformations {
						if cki.Name == ccki.Name {
							cki.Tags = append(cki.Tags, ccki.Tags...)
							cki.Logs = append(cki.Logs, ccki.Logs...)
							cki.Durations = append(cki.Durations, ccki.Durations...)
						}
					}
				}
			}
			trimmedMeanFraction := DefaultTrimmedMeanFraction
			cki.Duration = TrimmedMeanInt64Slice(cki.Durations, trimmedMeanFraction)
			cki.MeanFlops = GetMeanLogValue(cki, "flop_count_sp", trimmedMeanFraction)
			cki.MeanDramReadBytes = GetMeanLogValue(cki, "dram_read_bytes", trimmedMeanFraction)
			cki.MeanDramWriteBytes = GetMeanLogValue(cki, "dram_write_bytes", trimmedMeanFraction)
			layerGPUInfo.SummaryGPUKernelInformations[ii] = cki
		}
		layerGPUInfos = append(layerGPUInfos, layerGPUInfo)
	}

	summary = layerGPUInfos

	sort.Sort(summary)

	return summary, nil
}

func dummyPP() {
	// for importing pp
	pp.Println("dummy")
}

//

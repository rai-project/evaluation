package evaluation

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/rai-project/evaluation/writer"
	"github.com/rai-project/tracer"
	trace_tree "github.com/rai-project/tracer/convert"
	"github.com/spf13/cast"
	model "github.com/uber/jaeger/model/json"
)

type SummaryGPULayerInformation struct {
	SummaryLayerInformation `json:",inline"`
	SummaryGPUInformations  SummaryGPUInformations `json:"kernel_launch_information,omitempty"`
}

func (p SummaryGPULayerInformation) Len() int { return len(p.SummaryGPUInformations) }
func (p SummaryGPULayerInformation) Less(i, j int) bool {
	x := p.SummaryGPUInformations[i]
	y := p.SummaryGPUInformations[j]
	xDuration := TrimmedMeanInt64Slice(x.Durations, DefaultTrimmedMeanFraction)
	yDuration := TrimmedMeanInt64Slice(y.Durations, DefaultTrimmedMeanFraction)
	return xDuration > yDuration
}
func (p SummaryGPULayerInformation) Swap(i, j int) {
	p.SummaryGPUInformations[i], p.SummaryGPUInformations[j] = p.SummaryGPUInformations[j], p.SummaryGPUInformations[i]
}

type SummaryGPULayerInformations []SummaryGPULayerInformation

func (p SummaryGPULayerInformations) Len() int { return len(p) }
func (p SummaryGPULayerInformations) Less(i, j int) bool {
	return p[i].SummaryLayerInformation.Index < p[j].SummaryLayerInformation.Index
}
func (p SummaryGPULayerInformations) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (infos SummaryGPULayerInformations) Header(opts ...writer.Option) []string {
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

func (infos SummaryGPULayerInformations) Row(opts ...writer.Option) []string {
	return []string{}
}

func (infos SummaryGPULayerInformations) GetKernelLogKeys() []string {
	kernelLogs := []Metadata{}
	for _, info := range infos {
		for _, cudaKernelInformation := range info.SummaryGPUInformations {
			if len(cudaKernelInformation.Logs) == 0 {
				continue
			}
			kernelLogs = append(kernelLogs, cudaKernelInformation.Logs...)
		}
	}
	return getMetaDataKeys(kernelLogs)
}

// Rows ...
func (info SummaryGPULayerInformation) Rows(iopts ...writer.Option) [][]string {
	cudaKernelInfos := info.SummaryGPUInformations
	layerInfo := SummaryMeanLayerInformation(info.SummaryLayerInformation)
	layerInfoRow := layerInfo.Row(iopts...)

	opts := writer.NewOptions(iopts...)

	rows := [][]string{}

	kernelLogKeys := SummaryGPULayerInformations{info}.GetKernelLogKeys()

	isFilteredKernel := func(kernelInfo SummaryGPUInformation) bool {
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

	for _, cki := range cudaKernelInfos {
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

func (es Evaluations) SummaryGPULayerInformations(perfCol *PerformanceCollection) (SummaryGPULayerInformations, error) {
	summary := SummaryGPULayerInformations{}
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

	groupedLayerGPUInfos := make([][]SummaryGPULayerInformation, numGroups)
	for ii, grsp := range groupedSpans {
		if groupedLayerGPUInfos[ii] == nil {
			groupedLayerGPUInfos[ii] = []SummaryGPULayerInformation{}
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

			layerGPUInformation := SummaryGPULayerInformation{
				SummaryLayerInformation: layerInfo,
				SummaryGPUInformations:  []SummaryGPUInformation{},
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
				layerGPUInformation.SummaryGPUInformations = append(layerGPUInformation.SummaryGPUInformations, CUDALaunchSpantoGPUInformation(child))
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
				for infoIdx := range layerGPUInformation.SummaryGPUInformations {
					info := layerGPUInformation.SummaryGPUInformations[infoIdx]
					if info.CorrelationId != childCorrelationId {
						continue
					}
					// only record kernel duration when no gpu metrics are captured
					if len(info.Logs) == 0 {
						info.Durations = []int64{
							cast.ToInt64(child.Duration),
						}
					}
					layerGPUInformation.SummaryGPUInformations[infoIdx] = info
				}
			}
			groupedLayerGPUInfos[ii] = append(groupedLayerGPUInfos[ii], layerGPUInformation)
		}
	}

	layerGPUInfos := []SummaryGPULayerInformation{}
	for _, li := range groupedLayerGPUInfos[0] {
		layerGPUInfo := li
		for ii := range layerGPUInfo.SummaryGPUInformations {
			cki := layerGPUInfo.SummaryGPUInformations[ii]
			for _, lis := range groupedLayerGPUInfos[1:] {
				for _, lli := range lis {
					if lli.Name != li.Name || li.Index != li.Index {
						continue
					}
					for _, ccki := range lli.SummaryGPUInformations {
						if cki.Name == ccki.Name {
							cki.Tags = append(cki.Tags, ccki.Tags...)
							cki.Logs = append(cki.Logs, ccki.Logs...)
							cki.Durations = append(cki.Durations, ccki.Durations...)
						}
					}
				}
			}
			trimmedMeanFraction := DefaultTrimmedMeanFraction
			cki.MeanDuration = TrimmedMeanInt64Slice(cki.Durations, trimmedMeanFraction)
			cki.MeanFlops = GetMeanLogValue(cki, "flop_count_sp", trimmedMeanFraction)
			cki.MeanDramReadBytes = GetMeanLogValue(cki, "dram_read_bytes", trimmedMeanFraction)
			cki.MeanDramWriteBytes = GetMeanLogValue(cki, "dram_write_bytes", trimmedMeanFraction)
			layerGPUInfo.SummaryGPUInformations[ii] = cki
		}
		layerGPUInfos = append(layerGPUInfos, layerGPUInfo)
	}

	summary = layerGPUInfos

	return summary, nil
}

func dummyPP() {
	// for importing pp
	pp.Println("dummy")
}

//

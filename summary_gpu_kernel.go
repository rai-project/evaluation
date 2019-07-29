package evaluation

import (
	"encoding/json"
	"errors"
	"fmt"
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

//easyjson:json
type SummaryGPUKernelInformation struct {
	Name                  string     `json:"name,omitempty"`
	MangledName           string     `json:"mangled_name,omitempty"`
	Durations             []int64    `json:"durations,omitempty"`
	Tags                  []Metadata `json:"tags,omitempty"`
	Logs                  []Metadata `json:"logs,omitempty"`
	CorrelationId         int64      `json:"correlation_id,omitempty"`
	MeanDuration          float64    `json:"mean_duration,omitempty"`
	MeanFlops             float64    `json:"mean_flops,omitempty"`
	MeanDramReadBytes     float64    `json:"mean_dram_read_bytes,omitempty"`
	MeanDramWriteBytes    float64    `json:"mean_dram_write_bytes,omitempty"`
	MeanAchievedOccupancy float64    `json:"mean_achieved_occupancy,omitempty"`
	ArithmeticIntensity   float64    `json:"arithmetic_intensity,omitempty"`
	ArithmeticThroughput  float64    `json:"arithmetic_throughput,omitempty"`
	MemoryBound           bool       `json:"memory_bound,omitempty"`
}

type SummaryGPUKernelInformations []SummaryGPUKernelInformation

func (p SummaryGPUKernelInformations) Len() int { return len(p) }
func (p SummaryGPUKernelInformations) Less(i, j int) bool {
	return p[i].MeanDuration > p[j].MeanDuration
}
func (p SummaryGPUKernelInformations) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (info SummaryGPUKernelInformation) Header(opts ...writer.Option) []string {
	extraHeader := []string{
		"kernel_name",
		"kernel_duration (us)",
		"kernel_flops",
		"kernel_dram_read_bytes",
		"kernel_dram_write_bytes",
		"kernel_achieved_occupancy (%)",
		"kernel_arithmetic_intensity (flops/byte)",
		"kernel_arithmetic_throughput (GFlops)",
		"kernel_memory_bound",
		// "kernel_durations (us)",
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
		fmt.Sprintf("%.2f", info.MeanDuration),
		cast.ToString(info.MeanFlops),
		fmt.Sprintf("%.2f", info.MeanDramReadBytes),
		fmt.Sprintf("%.2f", info.MeanDramWriteBytes),
		fmt.Sprintf("%.2f", info.MeanAchievedOccupancy*100),
		fmt.Sprintf("%.2f", info.ArithmeticIntensity),
		fmt.Sprintf("%.2f", info.ArithmeticThroughput),
		cast.ToString(info.MemoryBound),
		// strings.Join(int64SliceToStringSlice(info.Durations), DefaultDimiter),
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

		kernelTags, err := json.Marshal(info.Tags)
		if err != nil {
			kernelTags = []byte{}
		}
		_ = kernelTags

		extra = append(extra, strings.Join(kernelLogs, DefaultDimiter))
	}
	return extra
}

//easyjson:json
type SummaryGPUKernelLayerInformation struct {
	SummaryLayerInformation      `json:",inline"`
	SummaryGPUKernelInformations SummaryGPUKernelInformations `json:"kernel_launch_information,omitempty"`
}

func (p SummaryGPUKernelLayerInformation) Len() int { return len(p.SummaryGPUKernelInformations) }
func (p SummaryGPUKernelLayerInformation) Less(i, j int) bool {
	x := p.SummaryGPUKernelInformations[i]
	y := p.SummaryGPUKernelInformations[j]
	return x.MeanDuration > y.MeanDuration
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
	extraHeader := SummaryGPUKernelInformation{}.Header()

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
		extra := cki.Row()
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
	kernelLogs := []float64{}
	for _, kernelLog := range info.Logs {
		for kernelLogKeyName, keryeLogValue := range kernelLog {
			if kernelLogKeyName == name {
				kernelLogs = append(kernelLogs, cast.ToFloat64(keryeLogValue))
			}
		}
	}
	return TrimmedMean(kernelLogs, trimmedMeanFraction)
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

	if len(es.GroupByBatchSize()) != 1 {
		return summary, errors.New("evaluations are not with the same batch size")
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
	for ii := range groupedLayerGPUInfos {
		if groupedLayerGPUInfos[ii] == nil {
			groupedLayerGPUInfos[ii] = []SummaryGPUKernelLayerInformation{}
		}
	}

	for ii, grsp := range groupedSpans {
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
				layerInfo = getLayerInfoFromLayerSpan(layerSpan)
				layerInfo.Durations = []int64{}
			} else {
				layerInfo = layerInfos.GetLayerInfoByName(layerSpan.OperationName)
			}

			layerGPUInformation := SummaryGPUKernelLayerInformation{
				SummaryLayerInformation:      layerInfo,
				SummaryGPUKernelInformations: []SummaryGPUKernelInformation{},
			}

			measureGPUMetrics := false

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
				info := CUDALaunchSpantoGPUInformation(child)
				if len(info.Logs) != 0 {
					measureGPUMetrics = true
				}
				layerGPUInformation.SummaryGPUKernelInformations = append(layerGPUInformation.SummaryGPUKernelInformations, info)
			}

			if !measureGPUMetrics {
				for _, ssp := range grsp {
					traceLevel, err := getTagValueAsString(ssp, "trace_level")
					if err != nil || traceLevel == "" {
						continue
					}
					if tracer.LevelFromName(traceLevel) != tracer.SYSTEM_LIBRARY_TRACE {
						continue
					}
					if strings.ToLower(ssp.OperationName) != "gpu_kernel" {
						continue
					}
					correlationId, err := getTagValueAsInt64(ssp, "correlation_id")
					if err != nil {
						log.WithError(err).Error("expecting cuda launch to have a correlation_id")
						continue
					}
					for infoIdx, _ := range layerGPUInformation.SummaryGPUKernelInformations {
						info := layerGPUInformation.SummaryGPUKernelInformations[infoIdx]
						if info.CorrelationId != correlationId {
							continue
						}
						info.Durations = []int64{
							cast.ToInt64(ssp.Duration),
						}
						layerGPUInformation.SummaryGPUKernelInformations[infoIdx] = info
					}
				}
			}

			groupedLayerGPUInfos[ii] = append(groupedLayerGPUInfos[ii], layerGPUInformation)
		}
	}

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
			cki.MeanDuration = TrimmedMeanInt64Slice(cki.Durations, trimmedMeanFraction)
			cki.MeanFlops = GetMeanLogValue(cki, "flop_count_sp", trimmedMeanFraction)
			cki.MeanDramReadBytes = GetMeanLogValue(cki, "dram_read_bytes", trimmedMeanFraction)
			cki.MeanDramWriteBytes = GetMeanLogValue(cki, "dram_write_bytes", trimmedMeanFraction)
			cki.MeanAchievedOccupancy = GetMeanLogValue(cki, "achieved_occupancy", trimmedMeanFraction)
			cki.ArithmeticIntensity = 0
			if (cki.MeanDramReadBytes + cki.MeanDramWriteBytes) != 0 {
				cki.ArithmeticIntensity = cki.MeanFlops / (cki.MeanDramReadBytes + cki.MeanDramWriteBytes)
			}
			cki.MemoryBound = false
			if cki.ArithmeticIntensity < layerGPUInfo.IdealArithmeticIntensity {
				cki.MemoryBound = true
			}
			cki.ArithmeticThroughput = cki.MeanFlops / cki.MeanDuration / float64(1000)
			layerGPUInfo.SummaryGPUKernelInformations[ii] = cki
		}
		summary = append(summary, layerGPUInfo)
	}

	sort.Sort(summary)

	return summary, nil
}

func dummyPP() {
	// for importing pp
	pp.Println("dummy")
}

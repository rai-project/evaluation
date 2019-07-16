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
	db "upper.io/db.v3"
)

type Metadata map[string]interface{}

type SummaryCUDAKernelInformation struct {
	Name          string     `json:"name,omitempty"`
	MangledName   string     `json:"mangled_name,omitempty"`
	Tags          []Metadata `json:"tags,omitempty"`
	Logs          []Metadata `json:"logs,omitempty"`
	Durations     []float64  `json:"durations,omitempty"`
	CorrelationId int64      `json:"correlation_id,omitempty"`
}

type SummaryCUDAKernelInformations []SummaryCUDAKernelInformation

type SummaryLayerCUDAKernelInformation struct {
	SummaryLayerInformation       `json:",inline"`
	SummaryCUDAKernelInformations SummaryCUDAKernelInformations `json:"kernel_launch_information,omitempty"`
}

func (p SummaryCUDAKernelInformations) Len() int { return len(p) }
func (p SummaryCUDAKernelInformations) Less(i, j int) bool {
	x := p[i]
	y := p[j]
	xDuration := TrimmedMean(x.Durations, DefaultTrimmedMeanFraction)
	yDuration := TrimmedMean(y.Durations, DefaultTrimmedMeanFraction)
	return xDuration > yDuration
}
func (p SummaryCUDAKernelInformations) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p SummaryLayerCUDAKernelInformation) Len() int { return len(p.SummaryCUDAKernelInformations) }
func (p SummaryLayerCUDAKernelInformation) Less(i, j int) bool {
	x := p.SummaryCUDAKernelInformations[i]
	y := p.SummaryCUDAKernelInformations[j]
	xDuration := TrimmedMean(x.Durations, DefaultTrimmedMeanFraction)
	yDuration := TrimmedMean(y.Durations, DefaultTrimmedMeanFraction)
	return xDuration > yDuration
}
func (p SummaryLayerCUDAKernelInformation) Swap(i, j int) {
	p.SummaryCUDAKernelInformations[i], p.SummaryCUDAKernelInformations[j] = p.SummaryCUDAKernelInformations[j], p.SummaryCUDAKernelInformations[i]
}

type SummaryLayerCUDAKernelInformations []SummaryLayerCUDAKernelInformation

func (p SummaryLayerCUDAKernelInformations) Len() int { return len(p) }
func (p SummaryLayerCUDAKernelInformations) Less(i, j int) bool {
	return p[i].SummaryLayerInformation.Index < p[j].SummaryLayerInformation.Index
}
func (p SummaryLayerCUDAKernelInformations) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (info SummaryLayerCUDAKernelInformations) Header(opts ...writer.Option) []string {
	extraHeader := []string{
		"kernel_name",
		"kernel_durations (us)",
	}
	if kernelLogKeys := getKernelLogKeys(info); len(kernelLogKeys) != 0 {
		extraHeader = append(extraHeader, kernelLogKeys...)
	}
	return append(SummaryLayerInformation{}.Header(opts...), extraHeader...)
}

func getKernelLogKeys(infos SummaryLayerCUDAKernelInformations) []string {
	kernelLogs := []Metadata{}
	for _, info := range infos {
		for _, cudaKernelInformation := range info.SummaryCUDAKernelInformations {
			if len(cudaKernelInformation.Logs) == 0 {
				continue
			}
			kernelLogs = append(kernelLogs, cudaKernelInformation.Logs...)
		}
	}
	return getMetaDataKeys(kernelLogs)
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

func getMetaDataValuesAsString(lg Metadata) []string {
	res := make([]string, len(lg))
	idx := 0
	for _, val := range lg {
		res[idx] = cast.ToString(val)
		idx += 1
	}
	return res
}

// Rows ...
func (info SummaryLayerCUDAKernelInformation) Rows(iopts ...writer.Option) [][]string {
	cudaKernelInfos := info.SummaryCUDAKernelInformations
	layerInfo := info.SummaryLayerInformation
	layerInfoRow := layerInfo.Row(iopts...)

	opts := writer.NewOptions(iopts...)

	rows := [][]string{}

	kernelLogKeys := getKernelLogKeys([]SummaryLayerCUDAKernelInformation{info})

	isFilteredKernel := func(kernelInfo SummaryCUDAKernelInformation) bool {
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
			strings.Join(float64SliceToStringSlice(cki.Durations), "\t"),
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
			extra = append(extra, strings.Join(kernelLogs, "\t"))
		}
		rows = append(rows, append(layerInfoRow, extra...))
	}
	return rows
}

func (k *SummaryCUDAKernelInformation) addLogs(spanLogs []model.Log) {
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

func (k *SummaryCUDAKernelInformation) addTags(spanTags []model.KeyValue) {
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

func toKernelInformation(span model.Span) SummaryCUDAKernelInformation {
	info := &SummaryCUDAKernelInformation{
		Name:          mustGetTagValueAsString(span, "kernel_name"),
		MangledName:   mustGetTagValueAsString(span, "name"),
		Tags:          []Metadata{},
		Logs:          []Metadata{},
		CorrelationId: mustGetTagValueAsInt64(span, "correlation_id"),
		Durations: []float64{
			cast.ToFloat64(span.Duration),
		},
	}
	info.addTags(span.Tags)
	info.addLogs(span.Logs)
	return *info
}

func (es Evaluations) GetSpansFromPerformanceCollection(perfCol *PerformanceCollection) (Spans, error) {
	spans := []model.Span{}
	for _, e := range es {
		foundPerfs, err := perfCol.Find(db.Cond{"_id": e.PerformanceID})
		if err != nil {
			return nil, err
		}
		if len(foundPerfs) != 1 {
			return nil, errors.New("expecting on performance output")
		}
		perf := foundPerfs[0]
		spans = append(spans, perf.Spans()...)
	}
	return spans, nil
}

func (es Evaluations) LayerCUDAKernelInformationSummary(perfCol *PerformanceCollection) (SummaryLayerCUDAKernelInformations, error) {
	summary := SummaryLayerCUDAKernelInformations{}
	if len(es) == 0 {
		return summary, errors.New("no evaluation is found in the database")
	}

	layerInfoSummary, err := es.SummaryLayerInformations(perfCol)
	if err != nil {
		log.WithError(err).Fatal("failed to get layerInformationSummary")
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

	groupedLayerCUDAKernelInfos := make([][]SummaryLayerCUDAKernelInformation, numGroups)
	for ii, grsp := range groupedSpans {
		if groupedLayerCUDAKernelInfos[ii] == nil {
			groupedLayerCUDAKernelInfos[ii] = []SummaryLayerCUDAKernelInformation{}
		}

		trace := model.Trace{
			TraceID: "0",
			Spans:   grsp,
		}
		tree, err := trace_tree.NewIntervalTree(trace)
		if err != nil {
			panic(err)
		}

		cnt := 0
		for _, sp := range grsp {
			traceLevel, err := getTagValueAsString(sp, "trace_level")
			if err != nil || traceLevel == "" {
				continue
			}
			if tracer.LevelFromName(traceLevel) != tracer.FRAMEWORK_TRACE {
				continue
			}

			cnt++

			layerSpan := trace_tree.ToInterval(sp)
			layerChildren := tree.ChildrenOf(layerSpan)
			layerCUDAKernelInformation := SummaryLayerCUDAKernelInformation{
				SummaryLayerInformation:       layerInfoSummary.GetLayerInfoByName(layerSpan.OperationName),
				SummaryCUDAKernelInformations: []SummaryCUDAKernelInformation{},
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

				child.Tags = append(child.Tags, model.KeyValue{
					Key:   "kernel_name",
					Type:  model.StringType,
					Value: demangleName(mustGetTagValueAsString(child, "name")),
				})
				layerCUDAKernelInformation.SummaryCUDAKernelInformations = append(layerCUDAKernelInformation.SummaryCUDAKernelInformations, toKernelInformation(child))
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
				childCorrelationId, err := getTagValueAsInt64(child, "correlation_id")
				if err != nil {
					log.WithError(err).Error("expecting cuda launch to have a correlation_id")
					continue
				}
				for infoIdx := range layerCUDAKernelInformation.SummaryCUDAKernelInformations {
					info := layerCUDAKernelInformation.SummaryCUDAKernelInformations[infoIdx]
					if info.CorrelationId != childCorrelationId {
						continue
					}
					info.addTags(child.Tags)
					info.addLogs(child.Logs)
					layerCUDAKernelInformation.SummaryCUDAKernelInformations[infoIdx] = info
				}
			}

			groupedLayerCUDAKernelInfos[ii] = append(groupedLayerCUDAKernelInfos[ii], layerCUDAKernelInformation)
		}
		pp.Println(cnt)

	}

	layerCUDAKernelInfos := []SummaryLayerCUDAKernelInformation{}
	for _, li := range groupedLayerCUDAKernelInfos[0] {
		layerCUDAKernelInfo := li
		for ii := range layerCUDAKernelInfo.SummaryCUDAKernelInformations {
			cki := layerCUDAKernelInfo.SummaryCUDAKernelInformations[ii]
			for _, lis := range groupedLayerCUDAKernelInfos[1:] {
				for _, lli := range lis {
					if lli.Name != li.Name || li.Index != li.Index {
						continue
					}
					for _, ccki := range lli.SummaryCUDAKernelInformations {
						if cki.Name == ccki.Name {
							cki.Tags = append(cki.Tags, ccki.Tags...)
							cki.Logs = append(cki.Logs, ccki.Logs...)
							cki.Durations = append(cki.Durations, ccki.Durations...)
						}
					}
				}
			}
			layerCUDAKernelInfo.SummaryCUDAKernelInformations[ii] = cki
		}
		layerCUDAKernelInfos = append(layerCUDAKernelInfos, layerCUDAKernelInfo)
	}

	summary = layerCUDAKernelInfos

	return summary, nil
}

func dummyPP() {
	// for import pp
	pp.Println("dummy")
}

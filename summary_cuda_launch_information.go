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

var summaryCUDAKernelInformationShowSummaryBase = false

type Metadata map[string]interface{}

type CUDAKernelInformation struct {
	Name          string     `json:"name,omitempty"`
	Tags          []Metadata `json:"tags,omitempty"`
	Logs          []Metadata `json:"logs,omitempty"`
	Durations     []float64  `json:"durations,omitempty"`
	CorrelationId int64      `json:"correlation_id,omitempty"`
}

type CUDAKernelInformations []CUDAKernelInformation

type LayerCUDAKernelInformation struct {
	LayerInformation       `json:",inline"`
	CUDAKernelInformations CUDAKernelInformations `json:"kernel_launch_information,omitempty"`
}

func (p CUDAKernelInformations) Len() int { return len(p) }
func (p CUDAKernelInformations) Less(i, j int) bool {
	x := p[i]
	y := p[j]
	xDuration := TrimmedMean(x.Durations, DefaultTrimmedMeanFraction)
	yDuration := TrimmedMean(y.Durations, DefaultTrimmedMeanFraction)
	return xDuration > yDuration
}
func (p CUDAKernelInformations) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p LayerCUDAKernelInformation) Len() int { return len(p.CUDAKernelInformations) }
func (p LayerCUDAKernelInformation) Less(i, j int) bool {
	x := p.CUDAKernelInformations[i]
	y := p.CUDAKernelInformations[j]
	xDuration := TrimmedMean(x.Durations, DefaultTrimmedMeanFraction)
	yDuration := TrimmedMean(y.Durations, DefaultTrimmedMeanFraction)
	return xDuration > yDuration
}
func (p LayerCUDAKernelInformation) Swap(i, j int) {
	p.CUDAKernelInformations[i], p.CUDAKernelInformations[j] = p.CUDAKernelInformations[j], p.CUDAKernelInformations[i]
}

type LayerCUDAKernelInformations []LayerCUDAKernelInformation

func (p LayerCUDAKernelInformations) Len() int { return len(p) }
func (p LayerCUDAKernelInformations) Less(i, j int) bool {
	return p[i].LayerInformation.Index < p[j].LayerInformation.Index
}
func (p LayerCUDAKernelInformations) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

type SummaryLayerCUDAKernelInformation struct {
	SummaryBase                 `json:",inline"`
	LayerCUDAKernelInformations LayerCUDAKernelInformations `json:"layer_informations,omitempty"`
}

func (info LayerCUDAKernelInformations) Header(opts ...writer.Option) []string {
	extraHeader := []string{
		"kernel_name",
		"kernel_durations (us)",
		// "kernel_tags",
		// "kernel_logs",
	}

	if kernelLogKeys := getKernelLogKeys(info); len(kernelLogKeys) != 0 {
		extraHeader = append(extraHeader, kernelLogKeys...)
	}
	return append(LayerInformation{}.Header(opts...), extraHeader...)
}

func (info0 LayerCUDAKernelInformation) Header(opts ...writer.Option) []string {
	info := LayerCUDAKernelInformations([]LayerCUDAKernelInformation{info0})
	return info.Header(opts...)
}

func getKernelLogKeys(infos LayerCUDAKernelInformations) []string {
	kernelLogs := []Metadata{}
	for _, info := range infos {
		for _, cudaKernelInformation := range info.CUDAKernelInformations {
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

func (info LayerCUDAKernelInformation) Rows(opts ...writer.Option) [][]string {
	cudaKernelInfos := info.CUDAKernelInformations
	layerInfo := info.LayerInformation
	layerInfoRow := layerInfo.Row(opts...)

	rows := [][]string{}

	kernelLogKeys := getKernelLogKeys([]LayerCUDAKernelInformation{info})

	for _, cki := range cudaKernelInfos {
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

func (LayerCUDAKernelInformations) Row(opts ...writer.Option) []string {
	panic("...")
	return nil
}

func (k *CUDAKernelInformation) addLogs(spanLogs []model.Log) {
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

func (k *CUDAKernelInformation) addTags(spanTags []model.KeyValue) {
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

func toKernelInformation(span model.Span) CUDAKernelInformation {
	info := &CUDAKernelInformation{
		Name:          mustGetTagValueAsString(span, "kernel_name"),
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

func (es Evaluations) LayerCUDAKernelInformationSummary(perfCol *PerformanceCollection) (SummaryLayerCUDAKernelInformation, error) {
	summary := SummaryLayerCUDAKernelInformation{}
	if len(es) == 0 {
		return summary, errors.New("no evaluation is found in the database")
	}

	summary = SummaryLayerCUDAKernelInformation{
		SummaryBase:                 es[0].summaryBase(),
		LayerCUDAKernelInformations: []LayerCUDAKernelInformation{},
	}

	spans, err := es.GetSpansFromPerformanceCollection(perfCol)
	if err != nil {
		return summary, err
	}
	if len(spans) == 0 {
		return summary, errors.New("no span is found for the evaluation")
	}

	predictSpans := spans.FilterByOperationName("c_predict")
	groupedSpans, err := getGroupedSpansFromSpans(predictSpans, spans)
	if err != nil {
		return summary, err
	}
	numGroups := len(groupedSpans)
	if numGroups == 0 {
		return summary, errors.New("no group of spans is found")
	}

	groupedLayerCUDAKernelInfos := make([][]LayerCUDAKernelInformation, numGroups)
	for ii, grsp := range groupedSpans {
		if groupedLayerCUDAKernelInfos[ii] == nil {
			groupedLayerCUDAKernelInfos[ii] = []LayerCUDAKernelInformation{}
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

			layerSpan := trace_tree.ToInterval(sp)
			layerChildren := tree.ChildrenOf(layerSpan)
			layerInfo, err := layerInformationSummary(es, []model.Span{predictSpans[ii], sp})
			if err != nil {
				log.WithError(err).Fatal("failed to get layerInformationSummary")
			}

			layerCUDAKernelInformation := LayerCUDAKernelInformation{
				LayerInformation:       layerInfo.LayerInformations[0],
				CUDAKernelInformations: []CUDAKernelInformation{},
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
				layerCUDAKernelInformation.CUDAKernelInformations = append(layerCUDAKernelInformation.CUDAKernelInformations, toKernelInformation(child))
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
				for infoIdx := range layerCUDAKernelInformation.CUDAKernelInformations {
					info := layerCUDAKernelInformation.CUDAKernelInformations[infoIdx]
					if info.CorrelationId != childCorrelationId {
						continue
					}
					info.addTags(child.Tags)
					info.addLogs(child.Logs)
					layerCUDAKernelInformation.CUDAKernelInformations[infoIdx] = info
				}
			}
			groupedLayerCUDAKernelInfos[ii] = append(groupedLayerCUDAKernelInfos[ii], layerCUDAKernelInformation)
		}
	}

	layerCUDAKernelInfos := []LayerCUDAKernelInformation{}
	for _, li := range groupedLayerCUDAKernelInfos[0] {
		layerCUDAKernelInfo := li
		for ii := range layerCUDAKernelInfo.CUDAKernelInformations {
			cki := layerCUDAKernelInfo.CUDAKernelInformations[ii]
			for _, lis := range groupedLayerCUDAKernelInfos[1:] {
				for _, lli := range lis {
					if lli.Name != li.Name || li.Index != li.Index {
						continue
					}
					for _, ccki := range lli.CUDAKernelInformations {
						if cki.Name == ccki.Name {
							cki.Tags = append(cki.Tags, ccki.Tags...)
							cki.Logs = append(cki.Logs, ccki.Logs...)
							cki.Durations = append(cki.Durations, ccki.Durations...)
						}
					}
				}
			}
			layerCUDAKernelInfo.CUDAKernelInformations[ii] = cki
		}
		layerCUDAKernelInfos = append(layerCUDAKernelInfos, layerCUDAKernelInfo)
	}

	summary.LayerCUDAKernelInformations = layerCUDAKernelInfos

	return summary, nil
}

func getGroupedCUDAKernelSpansFromSpans(predictSpans Spans, spans Spans) ([]Spans, error) {
	groupedSpans, err := getGroupedSpansFromSpans(predictSpans, spans)
	if err != nil {
		return nil, err
	}
	numPredictSpans := len(groupedSpans)

	groupedCUDAKernelSpans := make([]Spans, numPredictSpans)
	for ii, grsp := range groupedSpans {
		if len(grsp) == 0 {
			continue
		}

		groupedCUDAKernelSpans[ii] = Spans{}
		for _, sp := range grsp {
			traceLevel, err := getTagValueAsString(sp, "trace_level")
			if err != nil || traceLevel == "" {
				continue
			}
			if tracer.LevelFromName(traceLevel) != tracer.SYSTEM_LIBRARY_TRACE {
				continue
			}
			if strings.ToLower(sp.OperationName) != "gpu_kernel" {
				continue
			}
			groupedCUDAKernelSpans[ii] = append(groupedCUDAKernelSpans[ii], sp)
		}
	}

	return groupedCUDAKernelSpans, nil
}

func dummyPP() {
	// for import pp
	pp.Println("dummy")
}

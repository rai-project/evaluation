package evaluation

import (
	"encoding/json"
	"errors"
	"strings"

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

type LayerCUDAKernelInformations []LayerCUDAKernelInformation

type SummaryLayerCUDAKernelInformation struct {
	SummaryBase                 `json:",inline"`
	LayerCUDAKernelInformations LayerCUDAKernelInformations `json:"layer_informations,omitempty"`
}

func (LayerCUDAKernelInformation) Header() []string {
	extra := []string{
		"kernel_name",
		"kernel_durations (us)",
		// "kernel_tags",
		// "kernel_logs",
	}
	return append(LayerInformation{}.Header(), extra...)
}

func (info LayerCUDAKernelInformation) Rows() [][]string {
	cudaKernelInfos := info.CUDAKernelInformations
	layerInfo := info.LayerInformation

	rows := make([][]string, len(cudaKernelInfos))
	for ii, cki := range cudaKernelInfos {
		tags, err := json.Marshal(cki.Tags)
		if err != nil {
			tags = []byte{}
		}
		logs, err := json.Marshal(cki.Logs)
		if err != nil {
			logs = []byte{}
		}
		_ = tags
		_ = logs

		extra := []string{
			cki.Name,
			strings.Join(float64SliceToStringSlice(cki.Durations), "\t"),
			// string(tags),
			// string(logs),
		}
		row := append(layerInfo.Row(), extra...)
		rows[ii] = row
	}
	return rows
}

func (LayerCUDAKernelInformations) Header() []string {
	return LayerInformation{}.Header()
}

// func (p Performance) CUDAKernelInformationSummary(es Evaluations) ([]SummaryCUDAKernelInformation, error) {
// 	infoAcrossRuns := getGroupedCUDAKernelSpansFromSpans(es, p.Spans())
// 	numRuns := len(infoAcrossRuns)

// 	if numRuns == 0 {
// 		return nil, errors.New("no kernels found")
// 	}

// 	summaries := []SummaryCUDAKernelInformation{}

// 	getSummaryPosition := func(info SummaryCUDAKernelInformation) int {
// 		for ii, s := range summaries {
// 			if s.LayerInformation.Name == info.LayerInformation.Name && s.LayerInformation.Index == info.LayerInformation.Index {
// 				return ii
// 			}
// 		}
// 		return -1
// 	}

// 	for _, infoRun := range infoAcrossRuns {
// 		for _, layer := range infoRun {
// 			summaryPosition := getSummaryPosition(layer)
// 			if summaryPosition == -1 {
// 				summaryPosition = len(summaries)
// 				var summary SummaryCUDAKernelInformation
// 				deepcopy.Copy(&summary, layer)
// 				summaries = append(summaries, summary)
// 				continue
// 			}
// 			summary := summaries[summaryPosition]
// 			for ii := range summary.CUDAKernelInformations {
// 				summaryKernel := summary.CUDAKernelInformations[ii]
// 				for _, kernel := range layer.CUDAKernelInformations {
// 					if strings.ToLower(summaryKernel.Name) == strings.ToLower(kernel.Name) {
// 						summaryKernel.Logs = append(summaryKernel.Logs, kernel.Logs...)
// 						summaryKernel.Tags = append(summaryKernel.Tags, kernel.Tags...)
// 						summaryKernel.Durations = append(summaryKernel.Durations, kernel.Durations...)
// 					}
// 				}
// 				summary.CUDAKernelInformations[ii] = summaryKernel
// 			}
// 			summaries[summaryPosition] = summary
// 		}
// 	}
//
// 	return summaries, nil
// }

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
	k.Tags = append(k.Tags, tags)
}

func toKernelInformation(span model.Span) CUDAKernelInformation {
	logs := Metadata{}
	for _, v := range span.Logs {
		for _, f := range v.Fields {
			logs[f.Key] = f.Value
		}
	}
	metadata := Metadata{}
	for _, v := range span.Tags {
		metadata[v.Key] = v.Value
	}
	info := CUDAKernelInformation{
		Name:          mustGetTagValueAsString(span, "kernel_name"),
		Tags:          []Metadata{metadata},
		Logs:          []Metadata{logs},
		CorrelationId: mustGetTagValueAsInt64(span, "correlation_id"),
		Durations: []float64{
			cast.ToFloat64(span.Duration),
		},
	}
	return info
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
			layerInfo, err := layerInformationSummary(es, []model.Span{predictSpans[ii]})
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
				if strings.ToLower(child.OperationName) != "cuda_kernel" {
					continue
				}
				childCorrelationId, err := getTagValueAsInt64(child, "correlation_id")
				if err != nil {
					log.WithError(err).Error("expecting cuda launch to have a correlation_id")
					continue
				}
				for _, info := range layerCUDAKernelInformation.CUDAKernelInformations {
					if info.CorrelationId != childCorrelationId {
						continue
					}
					info.addTags(child.Tags)
					info.addLogs(child.Logs)
				}
			}

			groupedLayerCUDAKernelInfos[ii] = append(groupedLayerCUDAKernelInfos[ii], layerCUDAKernelInformation)
		}
	}

	layerCUDAKernelInfos := []LayerCUDAKernelInformation{}
	for _, ai := range groupedLayerCUDAKernelInfos[0] {
		layerCUDAKernelInfo := ai
		cudaKernelInfos := layerCUDAKernelInfo.CUDAKernelInformations
		for _, cki := range cudaKernelInfos {
			for _, lis := range groupedLayerCUDAKernelInfos[1:] {
				for _, li := range lis {
					if li.Name != ai.Name {
						continue
					}
					for _, ccki := range li.CUDAKernelInformations {
						if cki.Name == ccki.Name {
							cki.Tags = append(cki.Tags, ccki.Tags...)
							cki.Logs = append(cki.Logs, ccki.Logs...)
							cki.Durations = append(cki.Durations, ccki.Durations...)
						}
					}
				}
			}
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

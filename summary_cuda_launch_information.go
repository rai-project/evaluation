// +build ignore

package evaluation

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/getlantern/deepcopy"
	"github.com/rai-project/tracer"
	trace_tree "github.com/rai-project/tracer/convert"
	"github.com/spf13/cast"
	model "github.com/uber/jaeger/model/json"
	db "upper.io/db.v3"
)

var summaryCUDAKernelInformationShowSummaryBase = false

type Metadata map[string]interface{}

//easyjson:json
type CUDAKernelInformation struct {
	Name             string     `json:"name,omitempty"`
	Tags             []Metadata `json:"tags,omitempty"`
	Logs             []Metadata `json:"logs,omitempty"`
	Durations        []float64  `json:"durations,omitempty"`
	CorrelationId    int64      `json:"correlation_id,omitempty"`
	LayerInformation `json:",inline"`
}

type CUDAKernelInformations []CUDAKernelInformation

//easyjson:json
type SummaryCUDAKernelInformation struct {
	SummaryBase            `json:",inline"`
	CUDAKernelInformations CUDAKernelInformations `json:"kernel_launch_information,omitempty"`
}

type SummaryCUDAKernelInformations []SummaryCUDAKernelInformation

func (CUDAKernelInformation) Header() []string {
	return []string{
		"name",
		"tags",
		"logs",
		"durations (us)",
	}
}

func (info CUDAKernelInformation) Row() []string {
	tags, err := json.Marshal(info.Tags)
	if err != nil {
		tags = []byte{}
	}
	logs, err := json.Marshal(info.Logs)
	if err != nil {
		logs = []byte{}
	}
	_ = tags
	_ = logs
	return []string{
		info.Name,
		string(tags),
		string(logs),
		strings.Join(float64SliceToStringSlice(info.Durations), "\t"),
	}
}

func (CUDAKernelInformations) Header() []string {
	return CUDAKernelInformation{}.Header()
}

func (s CUDAKernelInformations) Rows() [][]string {
	rows := [][]string{}
	for _, e := range s {
		rows = append(rows, e.Row())
	}
	return rows
}

func (SummaryCUDAKernelInformation) Header() []string {
	extra := []string{
		"kernel_name",
		"kernel_duration",
		"layer_index",
		"layer_name",
		"layer_duration",
	}
	if summaryCUDAKernelInformationShowSummaryBase {
		return append(SummaryBase{}.Header(), extra...)
	}
	return extra
}

func (s SummaryCUDAKernelInformation) Row() []string {
	infos := []string{}
	for _, row := range s.CUDAKernelInformations.Rows() {
		infos = append(infos, strings.Join(row, ":"))
	}
	extra := []string{
		strings.Join(infos, ";"),
	}
	return append(s.SummaryBase.Row(), extra...)
}

// func (s SummaryCUDAKernelInformation) Rows() [][]string {
// 	infos := [][]string{}
// 	summaryRow := s.SummaryBase.Row()
// 	summaryRowLen := len(summaryRow)
// 	if !summaryCUDAKernelInformationShowSummaryBase {
// 		summaryRowLen = 0
// 	}
// 	rows := s.CUDAKernelInformations.Rows()
// 	infos = make([][]string, len(rows))
// 	for ii, row := range rows {
// 		infos[ii] = make([]string, summaryRowLen+len(row)+2)
// 		infos[ii][0] = cast.ToString(s.LayerInformation.Index)
// 		infos[ii][1] = s.LayerInformation.Name
// 		if summaryCUDAKernelInformationShowSummaryBase {
// 			for jj, elem := range summaryRow {
// 				infos[ii][jj+2] = elem
// 			}
// 		}
// 		for jj, elem := range row {
// 			infos[ii][jj+summaryRowLen+2] = elem
// 		}
// 	}
// 	return infos
// }

func (SummaryCUDAKernelInformations) Header() []string {
	return SummaryCUDAKernelInformation{}.Header()
}

func (s SummaryCUDAKernelInformations) Rows() [][]string {
	rows := [][]string{}
	for _, e := range s {
		rows = append(rows, e.Row())
	}
	return rows
}

// func (p Performance) CUDAKernelInformationSummary(es Evaluations) ([]SummaryCUDAKernelInformation, error) {
// 	infoAcrossRuns := getSpanKernelLaunchesFromSpans(es, p.Spans())
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

// func (e Evaluation) CUDAKernelInformationSummary(perfCol *PerformanceCollection) (SummaryCUDAKernelInformations, error) {
// 	perfs, err := perfCol.Find(db.Cond{"_id": e.PerformanceID})
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(perfs) != 1 {
// 		return nil, errors.New("expecting on performance output")
// 	}
// 	perf := perfs[0]
// 	return perf.CUDAKernelInformationSummary(e)
// }

func (es Evaluations) CUDAKernelInformationSummary(perfCol *PerformanceCollection) (SummaryCUDAKernelInformations, error) {
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

	infoAcrossRuns := getSpanKernelLaunchesFromSpans(es, spans)
	numRuns := len(infoAcrossRuns)

	if numRuns == 0 {
		return nil, errors.New("no kernels found")
	}

	summaries := []SummaryCUDAKernelInformation{}

	getSummaryPosition := func(info SummaryCUDAKernelInformation) int {
		for ii, s := range summaries {
			if s.LayerInformation.Name == info.LayerInformation.Name && s.LayerInformation.Index == info.LayerInformation.Index {
				return ii
			}
		}
		return -1
	}

	for _, infoRun := range infoAcrossRuns {
		for _, layer := range infoRun {
			summaryPosition := getSummaryPosition(layer)
			if summaryPosition == -1 {
				summaryPosition = len(summaries)
				var summary SummaryCUDAKernelInformation
				deepcopy.Copy(&summary, layer)
				summaries = append(summaries, summary)
				continue
			}
			summary := summaries[summaryPosition]
			for ii := range summary.CUDAKernelInformations {
				summaryKernel := summary.CUDAKernelInformations[ii]
				for _, kernel := range layer.CUDAKernelInformations {
					if strings.ToLower(summaryKernel.Name) == strings.ToLower(kernel.Name) {
						summaryKernel.Logs = append(summaryKernel.Logs, kernel.Logs...)
						summaryKernel.Tags = append(summaryKernel.Tags, kernel.Tags...)
						summaryKernel.Durations = append(summaryKernel.Durations, kernel.Durations...)
					}
				}
				summary.CUDAKernelInformations[ii] = summaryKernel
			}
			summaries[summaryPosition] = summary
		}
	}

	return summaries, nil
}

func getSpanKernelLaunchesFromSpans(es Evaluations, spans Spans) []SummaryCUDAKernelInformations {
	predictSpans := spans.FilterByOperationName("c_predict")
	numPredictSpans := len(predictSpans)
	groupedSpans := make([]Spans, numPredictSpans)
	for _, span := range spans {
		idx := predictSpanIndexOf(span, predictSpans)
		if idx == -1 {
			continue
		}
		var spanCopy model.Span
		deepcopy.Copy(&spanCopy, span)
		groupedSpans[idx] = append(groupedSpans[idx], spanCopy)
	}

	SummaryCUDAKernelInformations := make([]SummaryCUDAKernelInformations, numPredictSpans)
	for ii, grsp := range groupedSpans {
		trace := model.Trace{
			TraceID: "0",
			Spans:   grsp,
		}
		tree, err := trace_tree.NewIntervalTree(trace)
		if err != nil {
			panic(err)
		}

		SummaryCUDAKernelInformations[ii] = []SummaryCUDAKernelInformation{}
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
				log.WithError(err).Fatal("failed to get layerInformationSummary ")
			}
			SummaryCUDAKernelInformation := SummaryCUDAKernelInformation{
				SummaryBase:            es[0].summaryBase(),
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
				SummaryCUDAKernelInformation.CUDAKernelInformations = append(SummaryCUDAKernelInformation.CUDAKernelInformations, toKernelInformation(child))
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
				for _, kernel := range SummaryCUDAKernelInformation.CUDAKernelInformations {
					if kernel.CorrelationId != childCorrelationId {
						continue
					}
					kernel.addTags(child.Tags)
					kernel.addLogs(child.Logs)
				}
			}

			SummaryCUDAKernelInformations[ii] = append(SummaryCUDAKernelInformations[ii], SummaryCUDAKernelInformation)
		}
	}

	return SummaryCUDAKernelInformations
}

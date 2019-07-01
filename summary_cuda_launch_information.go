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

var summaryCUDALaunchInformationShowSummaryBase = false

type Metadata map[string]interface{}

//easyjson:json
type KernelLaunchInformation struct {
	Name          string     `json:"name,omitempty"`
	Tags          []Metadata `json:"tags,omitempty"`
	Logs          []Metadata `json:"logs,omitempty"`
	Durations     []float64  `json:"durations,omitempty"`
	CorrelationId int64      `json:"correlation_id,omitempty"`
}

type KernelLaunchInformations []KernelLaunchInformation

//easyjson:json
type SummaryCUDALaunchInformation struct {
	SummaryBase              `json:",inline"`
	LayerInformation         `json:",inline"`
	KernelLaunchInformations KernelLaunchInformations `json:"kernel_launch_information,omitempty"`
}

type SummaryCUDALaunchInformations []SummaryCUDALaunchInformation

func (KernelLaunchInformation) Header() []string {
	return []string{
		"name",
		// "tags",
		// "logs",
		"durations (us)",
	}
}

func (info KernelLaunchInformation) Row() []string {
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
		// string(tags),
		// string(logs),
		strings.Join(float64SliceToStringSlice(info.Durations), "\t"),
	}
}

func (KernelLaunchInformations) Header() []string {
	return KernelLaunchInformation{}.Header()
}

func (s KernelLaunchInformations) Rows() [][]string {
	rows := [][]string{}
	for _, e := range s {
		rows = append(rows, e.Row())
	}
	return rows
}

func (SummaryCUDALaunchInformation) Header() []string {
	extra := []string{
		"layer_index",
		"layer",
		"kernel_name",
		"duration",
	}
	if summaryCUDALaunchInformationShowSummaryBase {
		return append(SummaryBase{}.Header(), extra...)
	}
	return extra
}

func (s SummaryCUDALaunchInformation) Row() []string {
	infos := []string{}
	for _, row := range s.KernelLaunchInformations.Rows() {
		infos = append(infos, strings.Join(row, ":"))
	}
	extra := []string{
		strings.Join(infos, ";"),
	}
	return append(s.SummaryBase.Row(), extra...)
}

func (s SummaryCUDALaunchInformation) Rows() [][]string {
	infos := [][]string{}
	summaryRow := s.SummaryBase.Row()
	summaryRowLen := len(summaryRow)
	if !summaryCUDALaunchInformationShowSummaryBase {
		summaryRowLen = 0
	}
	rows := s.KernelLaunchInformations.Rows()
	infos = make([][]string, len(rows))
	for ii, row := range rows {
		infos[ii] = make([]string, summaryRowLen+len(row)+2)
		infos[ii][0] = s.LayerInformation.Name
		infos[ii][1] = cast.ToString(s.LayerInformation.Index)
		if summaryCUDALaunchInformationShowSummaryBase {
			for jj, elem := range summaryRow {
				infos[ii][jj+2] = elem
			}
		}
		for jj, elem := range row {
			infos[ii][jj+summaryRowLen+2] = elem
		}
	}
	return infos
}

func (SummaryCUDALaunchInformations) Header() []string {
	return SummaryCUDALaunchInformation{}.Header()
}

func (s SummaryCUDALaunchInformations) Rows() [][]string {
	rows := [][]string{}
	for _, e := range s {
		rows = append(rows, e.Row())
	}
	return rows
}

func (p Performance) CUDALaunchInformationSummary(e Evaluation) ([]SummaryCUDALaunchInformation, error) {
	es := []Evaluation{e}
	infoAcrossRuns := getSpanKernelLaunchesFromSpans(es, p.Spans())
	numRuns := len(infoAcrossRuns)

	if numRuns == 0 {
		return nil, errors.New("no kernels found")
	}

	summaries := []SummaryCUDALaunchInformation{}

	getSummaryPosition := func(info SummaryCUDALaunchInformation) int {
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
				var summary SummaryCUDALaunchInformation
				deepcopy.Copy(&summary, layer)
				summaries = append(summaries, summary)
				continue
			}
			summary := summaries[summaryPosition]
			for ii := range summary.KernelLaunchInformations {
				summaryKernel := summary.KernelLaunchInformations[ii]
				for _, kernel := range layer.KernelLaunchInformations {
					if strings.ToLower(summaryKernel.Name) == strings.ToLower(kernel.Name) {
						summaryKernel.Logs = append(summaryKernel.Logs, kernel.Logs...)
						summaryKernel.Tags = append(summaryKernel.Tags, kernel.Tags...)
						summaryKernel.Durations = append(summaryKernel.Durations, kernel.Durations...)
					}
				}
				summary.KernelLaunchInformations[ii] = summaryKernel
			}
			summaries[summaryPosition] = summary
		}
	}

	return summaries, nil
}

func (k *KernelLaunchInformation) addLogs(spanLogs []model.Log) {
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

func (k *KernelLaunchInformation) addTags(spanTags []model.KeyValue) {
	if k.Tags == nil {
		k.Tags = []Metadata{}
	}
	tags := Metadata{}
	for _, v := range spanTags {
		tags[v.Key] = v.Value
	}
	k.Tags = append(k.Tags, tags)
}

func toKernelInformation(span model.Span) KernelLaunchInformation {
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
	info := KernelLaunchInformation{
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

func (e Evaluation) CUDALaunchInformationSummary(perfCol *PerformanceCollection) (SummaryCUDALaunchInformations, error) {
	perfs, err := perfCol.Find(db.Cond{"_id": e.PerformanceID})
	if err != nil {
		return nil, err
	}
	if len(perfs) != 1 {
		return nil, errors.New("expecting on performance output")
	}
	perf := perfs[0]
	return perf.CUDALaunchInformationSummary(e)
}

func (es Evaluations) CUDALaunchInformationSummary(perfCol *PerformanceCollection) (SummaryCUDALaunchInformations, error) {
	res := []SummaryCUDALaunchInformation{}
	for _, e := range es {
		s, err := e.CUDALaunchInformationSummary(perfCol)
		if err != nil {
			log.WithError(err).
				WithField("framework_name", e.Framework.Name).
				WithField("model_name", e.Model.Name).
				Error("failed to get layer information summary")
			continue
		}
		if s == nil {
			continue
		}
		res = append(res, s...)
	}
	return res, nil
}

func getSpanKernelLaunchesFromSpans(es Evaluations, spans Spans) []SummaryCUDALaunchInformations {
	predictSpans := spans.FilterByOperationName("c_predict")
	groupedSpans := make([]Spans, len(predictSpans))
	for _, span := range spans {
		idx := predictIndexOf(span, predictSpans)
		if idx == -1 {
			continue
		}
		var spanCopy model.Span
		deepcopy.Copy(&spanCopy, span)
		groupedSpans[idx] = append(groupedSpans[idx], spanCopy)
	}

	summaryCUDALaunchInformations := make([]SummaryCUDALaunchInformations, len(predictSpans))
	for ii, grsp := range groupedSpans {
		trace := model.Trace{
			TraceID: "0",
			Spans:   grsp,
		}
		tree, err := trace_tree.NewIntervalTree(trace)
		if err != nil {
			panic(err)
		}

		summaryCUDALaunchInformations[ii] = []SummaryCUDALaunchInformation{}
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
			summaryCUDALaunchInformation := SummaryCUDALaunchInformation{
				SummaryBase:              es[0].summaryBase(),
				LayerInformation:         layerInfo.LayerInformations[0],
				KernelLaunchInformations: []KernelLaunchInformation{},
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
				summaryCUDALaunchInformation.KernelLaunchInformations = append(summaryCUDALaunchInformation.KernelLaunchInformations, toKernelInformation(child))
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
				for _, kernel := range summaryCUDALaunchInformation.KernelLaunchInformations {
					if kernel.CorrelationId != childCorrelationId {
						continue
					}
					kernel.addTags(child.Tags)
					kernel.addLogs(child.Logs)
				}
			}

			summaryCUDALaunchInformations[ii] = append(summaryCUDALaunchInformations[ii], summaryCUDALaunchInformation)

		}

	}

	return summaryCUDALaunchInformations
}

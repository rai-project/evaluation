package evaluation

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/k0kubun/pp"
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

		if false {
			pp.Println("logs = ", cki.Logs)
		}

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
							if ii < 10 {
								pp.Println(ccki.Logs)
							}
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

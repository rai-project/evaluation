package evaluation

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/rai-project/tracer"
	"github.com/spf13/cast"
	db "upper.io/db.v3"
)

type Metadata map[string]interface{}

type KernelLaunchInformation struct {
	Name      string     `json:"name,omitempty"`
	Tags      []Metadata `json:"tags,omitempty"`
	Logs      []Metadata `json:"logs,omitempty"`
	Durations []float64  `json:"durations,omitempty"`
}

type KernelLaunchInformations []KernelLaunchInformation

type SummaryCUDALaunchInformation struct {
	SummaryBase              `json:",inline"`
	KernelLaunchInformations KernelLaunchInformations `json:"kernel_launch_information,omitempty"`
}

type SummaryCUDALaunchInformations []SummaryCUDALaunchInformation

func (KernelLaunchInformation) Header() []string {
	return []string{
		"name",
		"tags",
		"logs",
		"durations",
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
	return []string{
		info.Name,
		string(tags),
		string(logs),
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
		"kernel_launch_information",
	}
	return append(SummaryBase{}.Header(), extra...)
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

func (p Performance) CUDALaunchInformationSummary(e Evaluation) (*SummaryCUDALaunchInformation, error) {
	sspans := getSpanLayersFromSpans(p.Spans())
	numSSpans := len(sspans)

	summary := &SummaryCUDALaunchInformation{
		SummaryBase:              e.summaryBase(),
		KernelLaunchInformations: KernelLaunchInformations{},
	}
	if numSSpans == 0 {
		return summary, nil
	}

	infosFull := make([][]KernelLaunchInformation, numSSpans)
	for ii, spans := range sspans {
		if infosFull[ii] == nil {
			infosFull[ii] = []KernelLaunchInformation{}
		}
		for _, span := range spans {
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
				Name: span.OperationName,
				Tags: []Metadata{metadata},
				Logs: []Metadata{logs},
				Durations: []float64{
					cast.ToFloat64(span.Duration),
				},
			}
			infosFull[ii] = append(infosFull[ii], info)
		}
	}

	infos := []KernelLaunchInformation{}
	for ii, span := range sspans[0] {
		durations := []float64{}
		tags := []Metadata{}
		logs := []Metadata{}
		for _, info := range infosFull {
			if len(info) <= ii {
				continue
			}
			tags = append(tags, info[ii].Tags...)
			logs = append(logs, info[ii].Logs...)
			durations = append(durations, info[ii].Durations...)
		}
		info := KernelLaunchInformation{
			Name:      span.OperationName,
			Tags:      tags,
			Logs:      logs,
			Durations: durations,
		}
		infos = append(infos, info)
	}

	summary.KernelLaunchInformations = infos
	return summary, nil
}

func (e Evaluation) CUDALaunchInformationSummary(perfCol *PerformanceCollection) (*SummaryCUDALaunchInformation, error) {
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
		res = append(res, *s)
	}
	return res, nil
}

func getSpanKernelLaunchesFromSpans(spans Spans) []Spans {
	predictSpans := spans.FilterByOperationName("Predict")
	groupedSpans := make([]Spans, len(predictSpans))
	for _, span := range spans {
		idx := predictIndexOf(span, predictSpans)
		if idx == -1 {
			continue
		}
		groupedSpans[idx] = append(groupedSpans[idx], span)
	}
	groupededSpans := make([]Spans, len(predictSpans))
	for ii, grp := range groupedSpans {
		groupededSpans[ii] = Spans{}
		if len(grp) == 0 {
			continue
		}
		predict := predictSpans[ii]
		traceLevel0, ok := spanTagValue(predict, "trace_level")
		if !ok {
			continue
		}
		traceLevel, ok := traceLevel0.(string)
		if !ok {
			continue
		}
		if traceLevel == "" {
			continue
		}
		if tracer.LevelFromName(traceLevel) < tracer.FRAMEWORK_TRACE {
			continue
		}
		r := groupededSpans[ii]
		for _, span := range grp {
			if strings.ToLower(span.OperationName) == "launch_kernel" {
				r = append(r, span)
			}
		}
		groupededSpans[ii] = r
	}
	return groupededSpans
}

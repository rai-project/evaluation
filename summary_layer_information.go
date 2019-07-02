package evaluation

import (
	json "encoding/json"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/k0kubun/pp"
	"github.com/rai-project/go-echarts/charts"
	"github.com/rai-project/tracer"
	"github.com/spf13/cast"
	model "github.com/uber/jaeger/model/json"
	db "upper.io/db.v3"
)

var (
	cntkLogMessageShown = false
)

//easyjson:json
type LayerInformation struct {
	Index     int       `json:"index,omitempty"`
	Name      string    `json:"name,omitempty"`
	Type      string    `json:"type,omitempty"`
	Durations []float64 `json:"durations,omitempty"`
}

type MeanLayerInformation struct {
	LayerInformation
}

//easyjson:json
type LayerInformations []LayerInformation

//easyjson:json
type MeanLayerInformations []MeanLayerInformation

//easyjson:json
type SummaryLayerInformation struct {
	SummaryBase       `json:",inline"`
	LayerInformations LayerInformations `json:"layer_informations,omitempty"`
}

func (LayerInformation) Header() []string {
	return []string{
		"layer_index",
		"layer_name",
		"layer_durations (us)",
	}
}

func (info LayerInformation) Row() []string {
	return []string{
		cast.ToString(info.Index),
		info.Name,
		strings.Join(float64SliceToStringSlice(info.Durations), ","),
	}
}

func (info MeanLayerInformation) Row() []string {
	return []string{
		cast.ToString(info.Index),
		info.Name,
		cast.ToString(TrimmedMean(info.Durations, 0)),
	}
}

func (LayerInformations) Header() []string {
	return LayerInformation{}.Header()
}

func (s LayerInformations) Rows() [][]string {
	rows := [][]string{}
	for _, e := range s {
		rows = append(rows, e.Row())
	}
	return rows
}

func layerInformationSummary(es Evaluations, spans Spans) (SummaryLayerInformation, error) {
	summary := SummaryLayerInformation{}
	if len(es) == 0 {
		return summary, errors.New("no evaluation is found in the database")
	}

	summary = SummaryLayerInformation{
		SummaryBase:       es[0].summaryBase(),
		LayerInformations: LayerInformations{},
	}

	predictSpans := spans.FilterByOperationName("c_predict")
	groupedLayerSpans, err := getGroupedLayerSpansFromSpans(predictSpans, spans)
	if err != nil {
		return summary, err
	}
	numGroups := len(groupedLayerSpans)
	if numGroups == 0 {
		return summary, errors.New("no group of spans is found")
	}

	groupedLayerInfos := make([][]LayerInformation, numGroups)
	for ii, spans := range groupedLayerSpans {
		if groupedLayerInfos[ii] == nil {
			groupedLayerInfos[ii] = []LayerInformation{}
		}
		for _, span := range spans {
			idx, err := getTagValueAsString(span, "layer_sequence_index")
			if err != nil || idx == "" {
				return summary, errors.New("cannot find tag layer_sequence_index")
			}
			layerInfo := LayerInformation{
				Index: cast.ToInt(idx),
				Name:  span.OperationName,
				Type:  getOpName(span),
				Durations: []float64{
					cast.ToFloat64(span.Duration),
				},
			}
			groupedLayerInfos[ii] = append(groupedLayerInfos[ii], layerInfo)
		}
	}

	layerInfos := []LayerInformation{}
	for ii, span := range groupedLayerSpans[0] {
		durations := []float64{}
		idx, err := getTagValueAsString(span, "layer_sequence_index")
		if err != nil || idx == "" {
			return summary, errors.New("cannot find tag layer_sequence_index")
		}
		for _, infos := range groupedLayerInfos {
			if len(infos) <= ii {
				continue
			}
			durationToAppend := []float64{}
			for _, info := range infos {
				if info.Index == cast.ToInt(idx) && info.Name == span.OperationName {
					durationToAppend = append(durationToAppend, info.Durations...)
				}
			}
			durations = append(durations, durationToAppend...)
		}

		layerInfos = append(layerInfos,
			LayerInformation{
				Index:     cast.ToInt(idx),
				Name:      span.OperationName,
				Type:      getOpName(span),
				Durations: durations,
			})
	}

	summary.LayerInformations = layerInfos

	return summary, nil
}

func (es Evaluations) LayerInformationSummary(perfCol *PerformanceCollection) (SummaryLayerInformation, error) {
	summary := SummaryLayerInformation{}
	spans := []model.Span{}
	for _, e := range es {
		foundPerfs, err := perfCol.Find(db.Cond{"_id": e.PerformanceID})
		if err != nil {
			return summary, err
		}
		if len(foundPerfs) != 1 {
			return summary, errors.New("no performance is found for the evaluation")
		}
		perf := foundPerfs[0]
		spans = append(spans, perf.Spans()...)
	}
	if len(spans) == 0 {
		return summary, errors.New("no span is found for the evaluation")
	}

	return layerInformationSummary(es, spans)
}

func sortByLayerIndex(spans Spans) {
	sort.Slice(spans, func(ii, jj int) bool {
		li, foundI := spanTagValue(spans[ii], "layer_sequence_index")
		if !foundI {
			return false
		}
		lj, foundJ := spanTagValue(spans[jj], "layer_sequence_index")
		if !foundJ {
			return true
		}

		return cast.ToInt64(li) < cast.ToInt64(lj)
	})
}

func getGroupedLayerSpansFromSpans(predictSpans Spans, spans Spans) ([]Spans, error) {
	groupedSpans, err := getGroupedSpansFromSpans(predictSpans, spans)
	if err != nil {
		return nil, err
	}
	numPredictSpans := len(groupedSpans)

	groupedLayerSpans := make([]Spans, numPredictSpans)
	for ii, grsp := range groupedSpans {
		if len(grsp) == 0 {
			continue
		}

		groupedLayerSpans[ii] = Spans{}
		for _, sp := range grsp {
			traceLevel, err := getTagValueAsString(sp, "trace_level")
			if err != nil || traceLevel == "" {
				continue
			}
			if tracer.LevelFromName(traceLevel) != tracer.FRAMEWORK_TRACE {
				continue
			}
			groupedLayerSpans[ii] = append(groupedLayerSpans[ii], sp)
		}

		sortByLayerIndex(groupedLayerSpans[ii])
	}

	return groupedLayerSpans, nil
}

func (o LayerInformations) BarPlot(title string) *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.TitleOpts{Title: title},
		charts.ToolboxOpts{Show: true},
	)
	bar = o.BarPlotAdd(bar)
	return bar
}

func (o LayerInformations) BarPlotAdd(bar *charts.Bar) *charts.Bar {
	timeUnit := time.Millisecond
	labels := []string{}
	for _, elem := range o {
		labels = append(labels, cast.ToString(elem.Index))
	}

	bar.AddXAxis(labels)
	durations := make([][]time.Duration, len(o[1].Durations))
	for ii := range o[1].Durations {
		durations[ii] = make([]time.Duration, len(o))
	}
	for ii, elem := range o {
		for jj, duration := range elem.Durations {
			durations[jj][ii] = time.Duration(duration) * time.Nanosecond / timeUnit
		}
	}
	for ii, duration := range durations {
		bar.AddYAxis(cast.ToString(ii), duration)
	}

	bar.SetSeriesOptions(charts.LabelTextOpts{Show: false})
	bar.SetGlobalOptions(
		charts.XAxisOpts{Name: "Layer Name"},
		charts.YAxisOpts{Name: "Latency(" + unitName(timeUnit) + ")"},
	)
	return bar
}

func (o LayerInformations) BoxPlot(title string) *charts.BoxPlot {
	box := charts.NewBoxPlot()
	box.SetGlobalOptions(
		charts.TitleOpts{Title: title},
		charts.ToolboxOpts{Show: true},
	)
	box = o.BoxPlotAdd(box)
	return box
}

func (o LayerInformations) BoxPlotAdd(box *charts.BoxPlot) *charts.BoxPlot {
	timeUnit := time.Nanosecond
	labels := []string{}
	for _, elem := range o {
		labels = append(labels, elem.Name)
	}
	box.AddXAxis(labels)

	durations := make([][]time.Duration, len(o))
	for ii, elem := range o {
		ts := make([]time.Duration, len(elem.Durations))
		for jj, t := range ts {
			ts[jj] = time.Duration(t) * timeUnit
		}
		durations[ii] = ts
	}
	box.AddYAxis("", durations)
	box.SetSeriesOptions(charts.LabelTextOpts{Show: false})
	box.SetGlobalOptions(
		charts.XAxisOpts{Name: "Layer Name"},
		charts.YAxisOpts{Name: "Latency(" + unitName(timeUnit) + ")"},
	)
	return box
}

func (o LayerInformations) Name() string {
	return "LayerInformations"
}

func (o LayerInformations) WriteBarPlot(path string) error {
	return writeBarPlot(o, path)
}

func (o LayerInformations) WriteBoxPlot(path string) error {
	return writeBoxPlot(o, path)
}

func (o LayerInformations) OpenBoxPlot() error {
	return openBoxPlot(o)
}

func (o LayerInformations) OpenBarPlot() error {
	return openBarPlot(o)
}

func (o MeanLayerInformations) BarPlot(title string) *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.TitleOpts{Title: title},
		charts.ToolboxOpts{Show: true},
	)
	bar = o.BarPlotAdd(bar)
	return bar
}

func (o MeanLayerInformations) BarPlotAdd(bar *charts.Bar) *charts.Bar {
	timeUnit := time.Nanosecond
	labels := []string{}
	for _, elem := range o {
		labels = append(labels, elem.Name)
	}
	bar.AddXAxis(labels)

	durations := make([]time.Duration, len(o))
	for ii, elem := range o {
		val := TrimmedMean(elem.Durations, 0)
		durations[ii] = time.Duration(val) * timeUnit
	}
	bar.AddYAxis("", durations)
	bar.SetSeriesOptions(charts.LabelTextOpts{Show: false})
	bar.SetGlobalOptions(
		charts.XAxisOpts{Name: "Layer Name"},
		charts.YAxisOpts{Name: "Latency(" + unitName(timeUnit) + ")"},
	)
	return bar
}

func (o MeanLayerInformations) BoxPlot(title string) *charts.BoxPlot {
	box := charts.NewBoxPlot()
	box.SetGlobalOptions(
		charts.TitleOpts{Title: title},
		charts.ToolboxOpts{Show: true},
	)
	box = o.BoxPlotAdd(box)
	return box
}

func (o MeanLayerInformations) BoxPlotAdd(box *charts.BoxPlot) *charts.BoxPlot {
	timeUnit := time.Nanosecond

	isPrivate := func(info MeanLayerInformation) bool {
		return strings.HasPrefix(info.Name, "_")
	}

	labels := []string{}
	for _, elem := range o {
		if isPrivate(elem) {
			continue
		}
		labels = append(labels, elem.Name)
	}
	box.AddXAxis(labels)

	durations := make([][]time.Duration, 0, len(o))
	for _, elem := range o {
		if isPrivate(elem) {
			continue
		}
		ts := make([]time.Duration, len(elem.Durations))
		for jj, t := range elem.Durations {
			ts[jj] = time.Duration(t) * timeUnit
		}
		durations = append(durations, prepareBoxplotData(ts))
	}
	if false {
		pp.Println(len(labels))
		pp.Println(len(durations))
		pp.Println(len(durations[0]))
	}
	box.AddYAxis("", durations)
	box.SetSeriesOptions(charts.LabelTextOpts{Show: false})
	jsLabelsBts, _ := json.Marshal(labels)
	jsFun := `function (name, index) {
    var labels = ` + strings.Replace(string(jsLabelsBts), `"`, "'", -1) + `;
    return labels.indexOf(name);
  }`
	box.SetGlobalOptions(
		charts.XAxisOpts{
			Name:      "Layer Name",
			Type:      "category",
			AxisLabel: charts.LabelTextOpts{Show: true, Rotate: 45, Formatter: charts.FuncOpts(jsFun)},
			SplitLine: charts.SplitLineOpts{Show: false},
			SplitArea: charts.SplitAreaOpts{Show: true},
		},
		charts.YAxisOpts{
			Name: "Latency(" + unitName(timeUnit) + ")",
			Type: "value",
			// NameRotate: 90,
			AxisLabel: charts.LabelTextOpts{Formatter: "{value}" + unitName(timeUnit)},
			SplitArea: charts.SplitAreaOpts{Show: true},
			Mix:       0,
		},
		charts.DataZoomOpts{
			Type:       "slider",
			XAxisIndex: []int{0},
			Start:      0,
			End:        float32(len(labels)),
		},
	)
	return box
}

func prepareBoxplotData(ds []time.Duration) []time.Duration {
	min := durationMin(ds)
	q1 := durationPercentile(ds, 25)
	q2 := durationPercentile(ds, 50)
	q3 := durationPercentile(ds, 75)
	max := durationMax(ds)
	return []time.Duration{min, q1, q2, q3, max}
}

func unitName(d time.Duration) string {
	return strings.TrimPrefix(d.String(), "1")
}

func (o MeanLayerInformations) Name() string {
	return "MeanLayerInformations"
}

func (o MeanLayerInformations) WriteBarPlot(path string) error {
	return writeBarPlot(o, path)
}

func (o MeanLayerInformations) WriteBoxPlot(path string) error {
	return writeBoxPlot(o, path)
}

func (o MeanLayerInformations) OpenBarPlot() error {
	return openBarPlot(o)
}

func (o MeanLayerInformations) OpenBoxPlot() error {
	return openBoxPlot(o)
}

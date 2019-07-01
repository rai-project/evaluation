package evaluation

import (
	json "encoding/json"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/getlantern/deepcopy"
	"github.com/k0kubun/pp"
	"github.com/rai-project/config"
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

//easyjson:json
type SummaryLayerInformations []SummaryLayerInformation

func (LayerInformation) Header() []string {
	return []string{
		"index",
		"name",
		"durations (us)",
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

func (SummaryLayerInformation) Header() []string {
	extra := []string{
		"layer_informations",
	}
	return append(SummaryBase{}.Header(), extra...)
}

func (s SummaryLayerInformation) Row() []string {
	infos := []string{}
	for _, row := range s.LayerInformations.Rows() {
		infos = append(infos, strings.Join(row, ":"))
	}
	extra := []string{
		strings.Join(infos, ";"),
	}
	return append(s.SummaryBase.Row(), extra...)
}

func (SummaryLayerInformations) Header() []string {
	return SummaryLayerInformation{}.Header()
}

func (s SummaryLayerInformations) Rows() [][]string {
	rows := [][]string{}
	for _, e := range s {
		rows = append(rows, e.Row())
	}
	return rows
}

func layerInformationSummary(es Evaluations, spans Spans) (*SummaryLayerInformation, error) {
	layerIndexIds := map[string]int{}
	for _, span := range spans {
		li, foundI := spanTagValue(span, "layer_sequence_index")
		if foundI {
			layerIndexIds[span.OperationName] = cast.ToInt(li)
		}
	}

	sspans := getSpanLayersFromSpans(spans)
	numSSpans := len(sspans)

	summary := &SummaryLayerInformation{
		SummaryBase:       es[0].summaryBase(),
		LayerInformations: LayerInformations{},
	}
	if numSSpans == 0 {
		return summary, nil
	}

	infosFull := make([][]LayerInformation, numSSpans)
	for ii, spans := range sspans {
		if infosFull[ii] == nil {
			infosFull[ii] = []LayerInformation{}
		}
		for _, span := range spans {
			info := LayerInformation{
				Index: layerIndexIds[span.OperationName],
				Name:  span.OperationName,
				Type:  getOpName(span),
				Durations: []float64{
					cast.ToFloat64(span.Duration),
				},
			}
			infosFull[ii] = append(infosFull[ii], info)
		}
	}

	infos := []LayerInformation{}
	for ii, span := range sspans[0] {
		durations := []float64{}
		for _, info := range infosFull {
			if len(info) <= ii {
				continue
			}
			durationToAppend := []float64{}
			for _, r := range info {
				if r.Name == span.OperationName {
					durationToAppend = append(durationToAppend, r.Durations...)
				}
			}
			durations = append(durations, durationToAppend...)
		}
		info := LayerInformation{
			Index:     layerIndexIds[span.OperationName],
			Name:      span.OperationName,
			Type:      getOpName(span),
			Durations: durations,
		}
		infos = append(infos, info)
	}

	summary.LayerInformations = infos
	return summary, nil
}

func getOpName(span model.Span) string {
	opName, err := getTagValueAsString(span, "op_name")
	if err != nil {
		return ""
	}
	return opName
}

func (p Performance) LayerInformationSummary(es Evaluations) (*SummaryLayerInformation, error) {
	return layerInformationSummary(es, p.Spans())
}

func (e Evaluation) LayerInformationSummary(perfCol *PerformanceCollection) (*SummaryLayerInformation, error) {
	perfs, err := perfCol.Find(db.Cond{"_id": e.PerformanceID})
	if err != nil {
		return nil, err
	}
	if len(perfs) != 1 {
		return nil, errors.New("expecting on performance output")
	}
	perf := perfs[0]
	return perf.LayerInformationSummary([]Evaluation{e})
}

func (es Evaluations) AcrossEvaluationLayerInformationSummary(perfCol *PerformanceCollection) (SummaryLayerInformations, error) {
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

	s, err := layerInformationSummary(es, spans)
	if err != nil {
		log.WithError(err).Error("failed to get layer information summary")
		return nil, err
	}
	if s == nil {
		return nil, errors.New("nil layer information summary")
	}
	return []SummaryLayerInformation{*s}, nil
}

func (es Evaluations) LayerInformationSummary(perfCol *PerformanceCollection) (SummaryLayerInformations, error) {
	res := []SummaryLayerInformation{}
	for _, e := range es {
		s, err := e.LayerInformationSummary(perfCol)
		if err != nil {
			log.WithError(err).Error("failed to get layer information summary")
			continue
		}
		if s == nil {
			continue
		}
		res = append(res, *s)
	}
	return res, nil
}

func spanIsCUPTI(span model.Span) bool {
	for _, tag := range span.Tags {
		key := strings.ToLower(tag.Key)
		switch key {
		case "cupti_domain", "cupti_callback_id":
			return true
		}
	}
	return false
}

func spanTagExists(span model.Span, key string) bool {
	for _, tag := range span.Tags {
		key0 := strings.ToLower(tag.Key)
		if key0 == key {
			return true
		}
	}
	return false
}

func spanTagValue(span model.Span, key string) (interface{}, bool) {
	for _, tag := range span.Tags {
		key0 := strings.ToLower(tag.Key)
		if key0 == key {
			return tag.Value, true
		}
	}
	return nil, false
}

func spanTagEquals(span model.Span, key string, value string) bool {
	for _, tag := range span.Tags {
		key0 := strings.ToLower(tag.Key)
		if key0 == key {
			e := strings.TrimSpace(strings.ToLower(cast.ToString(tag.Value)))
			return e == value
		}
	}
	return false
}

func spanComponentIs(span model.Span, name string) bool {
	for _, tag := range span.Tags {
		key := strings.ToLower(tag.Key)
		switch key {
		case "component":
			return strings.ToLower(cast.ToString(tag.Value)) == name
		}
	}
	return false
}

func selectTensorflowLayerSpans(spans Spans) Spans {
	res := []model.Span{}
	for _, span := range spans {
		if spanIsCUPTI(span) {
			continue
		}
		if !spanComponentIs(span, config.App.Name) {
			continue
		}
		if !spanTagExists(span, "thread_id") {
			continue
		}
		if !spanTagExists(span, "timeline_label") {
			continue
		}
		res = append(res, span)
	}
	return res
}

func selectMXNetLayerSpans(spans Spans) Spans {
	res := []model.Span{}
	for _, span := range spans {
		if spanIsCUPTI(span) {
			continue
		}
		if !spanComponentIs(span, config.App.Name) {
			continue
		}
		if !spanTagExists(span, "thread_id") {
			continue
		}
		if !spanTagExists(span, "process_id") {
			continue
		}
		res = append(res, span)
	}
	return res
}
func selectCaffeLayerSpans(spans Spans) Spans {
	return selectCaffe2LayerSpans(spans)
}
func selectCaffe2LayerSpans(spans Spans) Spans {
	res := []model.Span{}
	for _, span := range spans {
		if spanIsCUPTI(span) {
			continue
		}
		if !spanComponentIs(span, config.App.Name) {
			continue
		}
		if !spanTagExists(span, "metadata") {
			continue
		}
		if !spanTagExists(span, "thread_id") {
			continue
		}
		if !spanTagEquals(span, "process_id", "0") {
			continue
		}
		res = append(res, span)
	}
	return res
}

func selectCNTKLayerSpans(spans Spans) Spans {
	if cntkLogMessageShown {
		return Spans{}
	}
	cntkLogMessageShown = true
	log.WithField("function", "selectCNTKLayerSpans").Error("layer information is not currently supported by cntk")
	return Spans{}
}
func selectTensorRTLayerSpans(spans Spans) Spans {
	return selectCaffe2LayerSpans(spans)
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

func getSpanLayersFromSpans(spans Spans) []Spans {
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
	groupedLayerSpans := make([]Spans, len(predictSpans))
	for ii, grsp := range groupedSpans {
		groupedLayerSpans[ii] = Spans{}
		if len(grsp) == 0 {
			continue
		}
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

	return groupedLayerSpans
}

func frameworkNameOfSpan(predictSpan model.Span) string {
	tagName := "framework_name"
	for _, tag := range predictSpan.Tags {
		if tag.Key == tagName {
			return cast.ToString(tag.Value)
		}
	}
	return ""
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
		labels = append(labels, elem.Name)
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
		charts.YAxisOpts{Name: "Latency(" + timeUnit.String() + ")"},
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
	timeUnit := time.Microsecond
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
		charts.YAxisOpts{Name: "Latency(" + timeUnit.String() + ")"},
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
	timeUnit := time.Microsecond
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
	timeUnit := time.Microsecond

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
    var labels = ` + strings.ReplaceAll(string(jsLabelsBts), `"`, "'") + `;
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

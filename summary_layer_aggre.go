package evaluation

import (
	"math"

	"github.com/rai-project/evaluation/writer"
	"github.com/rai-project/go-echarts/charts"
	"github.com/spf13/cast"
)

//easyjson:json
type SummaryLayerAggreInformation struct {
	SummaryModelInformation `json:",inline"`
	Type                    string  `json:"type,omitempty"`
	Occurence               int     `json:"occurrence,omitempty"`
	OccurencePercentage     float64 `json:"occurrence_percentage,omitempty"`
	Duration                float64 `json:"duration,omitempty"`
	DurationPercentage      float64 `json:"duration_percentage,omitempty"`
	Memory                  int64   `json:"memory,omitempty"`
	MemoryPercentage        float64 `json:"memory_percentage,omitempty"`
}

//easyjson:json
type SummaryLayerAggreInformations []SummaryLayerAggreInformation

//easyjson:json
type SummaryLayerDruationInformation SummaryLayerAggreInformation

//easyjson:json
type SummaryLayerOccurenceInformation SummaryLayerAggreInformation

//easyjson:json
type SummaryLayerDruationInformations SummaryLayerAggreInformations

//easyjson:json
type SummaryLayerOccurrenceInformations SummaryLayerAggreInformations

func (SummaryLayerAggreInformation) Header(iopts ...writer.Option) []string {
	extra := []string{
		"type",
		"occurrences",
		"occurrence percentage (%)",
		"duration (us)",
		"duration percentage (%)",
	}
	opts := writer.NewOptions(iopts...)
	if opts.ShowSummaryBase {
		return append(SummaryBase{}.Header(iopts...), extra...)
	}
	return extra
}

func (s SummaryLayerAggreInformation) Row(iopts ...writer.Option) []string {
	extra := []string{
		s.Type,
		cast.ToString(s.Occurence),
		cast.ToString(s.OccurencePercentage),
		cast.ToString(s.Duration),
		cast.ToString(s.DurationPercentage),
	}
	opts := writer.NewOptions(iopts...)
	if opts.ShowSummaryBase {
		return append(s.SummaryBase.Header(iopts...), extra...)
	}
	return extra
}

func (es Evaluations) SummaryLayerAggreInformations(perfCol *PerformanceCollection) (SummaryLayerAggreInformations, error) {
	summary := SummaryLayerAggreInformations{}
	layerInfos, err := es.SummaryLayerInformations(perfCol)
	if err != nil {
		return summary, err
	}

	modelInfos, err := (es.SummaryModelInformations(perfCol))
	modelInfo := modelInfos[0]
	if err != nil {
		modelInfo = SummaryModelInformation{}
	}

	exsistedLayers := make(map[string]SummaryLayerAggreInformation)
	totalOcurrences := 0
	totalDuration := float64(0)
	for _, info := range layerInfos {
		layerType := info.Type
		duration := TrimmedMeanInt64Slice(info.Durations, DefaultTrimmedMeanFraction)
		v, ok := exsistedLayers[layerType]
		if !ok {
			exsistedLayers[layerType] = SummaryLayerAggreInformation{
				SummaryModelInformation: modelInfo,
				Type:                    layerType,
				Occurence:               1,
				Duration:                duration,
			}
		} else {
			v.Occurence += 1
			v.Duration += duration
			exsistedLayers[layerType] = v
		}
		totalOcurrences += 1
		totalDuration += duration
	}

	for _, info := range exsistedLayers {
		info.DurationPercentage = math.Round(100*100*float64(info.Duration)/float64(totalDuration)) / 100
		info.OccurencePercentage = math.Round(100*100*float64(info.Occurence)/float64(totalOcurrences)) / 100
		summary = append(summary, info)
	}

	return summary, nil
}

func (o SummaryLayerDruationInformations) PlotName() string {
	if len(o) == 0 {
		return ""
	}
	return o[0].ModelName + " Batch Size = " + cast.ToString(o[0].BatchSize) + " Layer Duration Percentage"
}

func (o SummaryLayerDruationInformations) PiePlot(title string) *charts.Pie {
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.TitleOpts{Title: title},
	)
	pie = o.PiePlotAdd(pie)
	return pie
}

func (o SummaryLayerOccurrenceInformations) PlotName() string {
	if len(o) == 0 {
		return ""
	}
	return o[0].ModelName + " Layer Occurrence Percentage"
}

func (o SummaryLayerOccurrenceInformations) PiePlot(title string) *charts.Pie {
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.TitleOpts{Title: title},
	)
	pie = o.PiePlotAdd(pie)
	return pie
}

type LayerAggregatedInformationSelector func(elem SummaryLayerAggreInformation) interface{}

func (o SummaryLayerAggreInformations) piePlotAdd(pie *charts.Pie, elemSelector LayerAggregatedInformationSelector) *charts.Pie {
	labels := []string{}
	data := make(map[string]interface{})
	for _, elem := range o {
		label := cast.ToString(elem.Type)
		data[label] = elemSelector(elem)
		labels = append(labels, label)
	}
	pie.AddSorted("", data, charts.LabelTextOpts{Show: true, Formatter: "{b}: {c}"})
	return pie
}

func (o SummaryLayerDruationInformations) PiePlotAdd(pie *charts.Pie) *charts.Pie {
	return SummaryLayerAggreInformations(o).piePlotAdd(pie, func(elem SummaryLayerAggreInformation) interface{} {
		return elem.DurationPercentage
	})
}
func (o SummaryLayerOccurrenceInformations) PiePlotAdd(pie *charts.Pie) *charts.Pie {
	return SummaryLayerAggreInformations(o).piePlotAdd(pie, func(elem SummaryLayerAggreInformation) interface{} {
		return elem.OccurencePercentage
	})
}

func (o SummaryLayerOccurrenceInformations) WritePiePlot(path string) error {
	return writePiePlot(o, path)
}

func (o SummaryLayerOccurrenceInformations) OpenPiePlot() error {
	return openPiePlot(o)
}

func (o SummaryLayerDruationInformations) WritePiePlot(path string) error {
	return writePiePlot(o, path)
}

func (o SummaryLayerDruationInformations) OpenPiePlot() error {
	return openPiePlot(o)
}

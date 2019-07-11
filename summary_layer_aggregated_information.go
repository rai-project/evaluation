package evaluation

import (
	"math"

	"github.com/rai-project/evaluation/writer"
	"github.com/rai-project/go-echarts/charts"
	"github.com/spf13/cast"
)

//easyjson:json
type SummaryLayerAggregatedInformation struct {
	SummaryBase         `json:",inline"`
	Type                string  `json:"type,omitempty"`
	Occurence           int     `json:"occurrence,omitempty"`
	OccurencePercentage float64 `json:"occurrence_percentage,omitempty"`
	Duration            float64 `json:"duration,omitempty"`
	DurationPercentage  float64 `json:"duration_percentage,omitempty"`
}

//easyjson:json
type SummaryLayerAggregatedInformations []SummaryLayerAggregatedInformation

//easyjson:json
type SummaryLayerDruationInformation SummaryLayerAggregatedInformation

//easyjson:json
type SummaryLayerOccurenceInformation SummaryLayerAggregatedInformation

//easyjson:json
type SummaryLayerDruationInformations SummaryLayerAggregatedInformations

//easyjson:json
type SummaryLayerOccurrenceInformations SummaryLayerAggregatedInformations

func (SummaryLayerAggregatedInformation) Header(iopts ...writer.Option) []string {
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

func (s SummaryLayerAggregatedInformation) Row(iopts ...writer.Option) []string {
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

func (es Evaluations) SummaryLayerAggregatedInformation(perfCol *PerformanceCollection) (SummaryLayerAggregatedInformations, error) {
	summary := SummaryLayerAggregatedInformations{}
	layerInfos, err := es.SummaryLayerInformations(perfCol)
	if err != nil {
		return summary, err
	}
	summaryBase := layerInfos[0].SummaryBase

	exsistedLayers := make(map[string]SummaryLayerAggregatedInformation)
	totalOcurrences := 0
	totalDuration := float64(0)
	for _, info := range layerInfos {
		layerType := info.Type
		duration := TrimmedMean(info.Durations, DefaultTrimmedMeanFraction)
		v, ok := exsistedLayers[layerType]
		if !ok {
			exsistedLayers[layerType] = SummaryLayerAggregatedInformation{
				SummaryBase: summaryBase,
				Type:        layerType,
				Occurence:   1,
				Duration:    duration,
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
		info.DurationPercentage = math.Round(100*100*info.Duration/totalDuration) / 100
		info.OccurencePercentage = math.Round(100*100*float64(info.Occurence)/float64(totalOcurrences)) / 100
		summary = append(summary, info)
	}

	return summary, nil
}

func (o SummaryLayerDruationInformations) Name() string {
	if len(o) == 0 {
		return ""
	}
	return o[0].ModelName + " Layer Duration Percentage"
}

func (o SummaryLayerDruationInformations) PiePlot(title string) *charts.Pie {
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.TitleOpts{Title: title},
	)
	pie = o.PiePlotAdd(pie)
	return pie
}

func (o SummaryLayerOccurrenceInformations) Name() string {
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

type LayerAggregatedInformationSelector func(elem SummaryLayerAggregatedInformation) interface{}

func (o SummaryLayerAggregatedInformations) piePlotAdd(pie *charts.Pie, elemSelector LayerAggregatedInformationSelector) *charts.Pie {
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
	return SummaryLayerAggregatedInformations(o).piePlotAdd(pie, func(elem SummaryLayerAggregatedInformation) interface{} {
		return elem.DurationPercentage
	})
}
func (o SummaryLayerOccurrenceInformations) PiePlotAdd(pie *charts.Pie) *charts.Pie {
	return SummaryLayerAggregatedInformations(o).piePlotAdd(pie, func(elem SummaryLayerAggregatedInformation) interface{} {
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

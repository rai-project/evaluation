package evaluation

import (
	"github.com/rai-project/evaluation/writer"
	"github.com/rai-project/go-echarts/charts"
	"github.com/spf13/cast"
)

type LayerAggregatedInformation struct {
	Type                string  `json:"type,omitempty"`
	Occurences          int     `json:"occurences,omitempty"`
	OccurencePercentage float32 `json:"occurence_percentage,omitempty"`
	Duration            float64 `json:"duration,omitempty"`
	DurationPercentage  float32 `json:"duration_percentage,omitempty"`
}

type LayerAggregatedInformations []LayerAggregatedInformation

//easyjson:json
type SummaryLayerAggregatedInformation struct {
	SummaryBase                 `json:",inline"`
	LayerAggregatedInformations LayerAggregatedInformations `json:"layer_aggregated_informations,omitempty"`
}

type SummaryLayerAggregatedInformationOccurrences struct {
	SummaryLayerAggregatedInformation
}

func (LayerAggregatedInformation) Header(opts ...writer.Option) []string {
	return []string{
		"type",
		"occurences",
		"occurence percentage (%)",
		"duration (us)",
		"duration percentage (%)",
	}
}

func (info LayerAggregatedInformation) Row(opts ...writer.Option) []string {
	return []string{
		info.Type,
		cast.ToString(info.Occurences),
		cast.ToString(info.OccurencePercentage),
		cast.ToString(info.Duration),
		cast.ToString(info.DurationPercentage),
	}
}

func (es Evaluations) LayerAggregatedInformationSummary(perfCol *PerformanceCollection) (SummaryLayerAggregatedInformation, error) {
	summary := SummaryLayerAggregatedInformation{}
	layerinfoSum, err := es.LayerInformationSummary(perfCol)
	if err != nil {
		return summary, err
	}
	layerInfos := layerinfoSum.LayerInformations

	summary = SummaryLayerAggregatedInformation{
		SummaryBase:                 es[0].summaryBase(),
		LayerAggregatedInformations: LayerAggregatedInformations{},
	}

	exsistedLayers := make(map[string]LayerAggregatedInformation)
	totalOcurrences := 0
	totalDuration := float64(0)
	for _, info := range layerInfos {
		layerType := info.Type
		duration := TrimmedMean(info.Durations, DefaultTrimmedMeanFraction)
		v, ok := exsistedLayers[layerType]
		if !ok {
			exsistedLayers[layerType] = LayerAggregatedInformation{
				Type:       layerType,
				Occurences: 1,
				Duration:   duration,
			}
		} else {
			v.Occurences += 1
			v.Duration += duration
			exsistedLayers[layerType] = v
		}
		totalOcurrences += 1
		totalDuration += duration
	}

	layerAggreInfos := []LayerAggregatedInformation{}
	for _, info := range exsistedLayers {
		info.DurationPercentage = 100 * float32(info.Duration/totalDuration)
		info.OccurencePercentage = 100 * float32(info.Occurences) / float32(totalOcurrences)
		layerAggreInfos = append(layerAggreInfos, info)
	}

	summary.LayerAggregatedInformations = layerAggreInfos

	return summary, nil
}

func (o SummaryLayerAggregatedInformation) Name() string {
	return o.ModelName + " Layer Type Composition"
}

func (o SummaryLayerAggregatedInformation) PiePlot(title string) *charts.Pie {
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.TitleOpts{Title: title},
	)
	pie = o.PiePlotAdd(pie)
	return pie
}

func (o SummaryLayerAggregatedInformationOccurrences) Name() string {
	return o.ModelName + " Layer Type Occurrences"
}

func (o SummaryLayerAggregatedInformationOccurrences) PiePlot(title string) *charts.Pie {
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.TitleOpts{Title: title},
	)
	pie = o.PiePlotAdd(pie)
	return pie
}

type LayerAggregatedInformationSelector func(elem LayerAggregatedInformation) interface{}

func (o SummaryLayerAggregatedInformation) piePlotAdd(pie *charts.Pie, elemSelector LayerAggregatedInformationSelector) *charts.Pie {
	labels := []string{}
	data := make(map[string]interface{})
	for _, elem := range o.LayerAggregatedInformations {
		label := cast.ToString(elem.Type)
		data[label] = elemSelector(elem)
		labels = append(labels, label)

	}
	pie.AddSorted("", data, charts.LabelTextOpts{Show: true})
	return pie
}

func (o SummaryLayerAggregatedInformation) PiePlotAdd(pie *charts.Pie) *charts.Pie {
	return o.piePlotAdd(pie, func(elem LayerAggregatedInformation) interface{} {
		return elem.Duration
	})
}

func (o SummaryLayerAggregatedInformationOccurrences) PiePlotAdd(pie *charts.Pie) *charts.Pie {
	return o.piePlotAdd(pie, func(elem LayerAggregatedInformation) interface{} {
		return elem.Occurences
	})
}

func (o SummaryLayerAggregatedInformationOccurrences) WritePiePlot(path string) error {
	return writePiePlot(o, path)
}

func (o SummaryLayerAggregatedInformationOccurrences) OpenPiePlot() error {
	return openPiePlot(o)
}

func (o SummaryLayerAggregatedInformation) WritePiePlot(path string) error {
	return writePiePlot(o, path)
}

func (o SummaryLayerAggregatedInformation) OpenPiePlot() error {
	return openPiePlot(o)
}

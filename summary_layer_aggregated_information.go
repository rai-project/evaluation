package evaluation

import (
	"github.com/rai-project/evaluation/writer"
	"github.com/rai-project/go-echarts/charts"
	"github.com/spf13/cast"
)

type LayerAggregatedInformation struct {
	Type                string  `json:"type,omitempty"`
	Occurence           int     `json:"occurrence,omitempty"`
	OccurencePercentage float32 `json:"occurrence_percentage,omitempty"`
	Duration            float64 `json:"duration,omitempty"`
	DurationPercentage  float32 `json:"duration_percentage,omitempty"`
}

type LayerAggregatedInformations []LayerAggregatedInformation

//easyjson:json
type SummaryLayerDruationInformation struct {
	SummaryBase                 `json:",inline"`
	LayerAggregatedInformations LayerAggregatedInformations `json:"layer_aggregated_informations,omitempty"`
}

type SummaryLayerOccurenceInformation struct {
	SummaryLayerDruationInformation `json:",inline"`
}

func (LayerAggregatedInformation) Header(opts ...writer.Option) []string {
	return []string{
		"type",
		"occurrences",
		"occurrence percentage (%)",
		"duration (us)",
		"duration percentage (%)",
	}
}

func (info LayerAggregatedInformation) Row(opts ...writer.Option) []string {
	return []string{
		info.Type,
		cast.ToString(info.Occurence),
		cast.ToString(info.OccurencePercentage),
		cast.ToString(info.Duration),
		cast.ToString(info.DurationPercentage),
	}
}

func (es Evaluations) LayerAggregatedInformationSummary(perfCol *PerformanceCollection) (SummaryLayerDruationInformation, error) {
	summary := SummaryLayerDruationInformation{}
	layerinfoSum, err := es.LayerInformationSummary(perfCol)
	if err != nil {
		return summary, err
	}
	layerInfos := layerinfoSum.LayerInformations

	summary = SummaryLayerDruationInformation{
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
				Type:      layerType,
				Occurence: 1,
				Duration:  duration,
			}
		} else {
			v.Occurence += 1
			v.Duration += duration
			exsistedLayers[layerType] = v
		}
		totalOcurrences += 1
		totalDuration += duration
	}

	layerAggreInfos := []LayerAggregatedInformation{}
	for _, info := range exsistedLayers {
		info.DurationPercentage = 100 * float32(info.Duration/totalDuration)
		info.OccurencePercentage = 100 * float32(info.Occurence) / float32(totalOcurrences)
		layerAggreInfos = append(layerAggreInfos, info)
	}

	summary.LayerAggregatedInformations = layerAggreInfos

	return summary, nil
}

func (o SummaryLayerDruationInformation) Name() string {
	return o.ModelName + " Layer Duration"
}

func (o SummaryLayerDruationInformation) PiePlot(title string) *charts.Pie {
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.TitleOpts{Title: title},
	)
	pie = o.PiePlotAdd(pie)
	return pie
}

func (o SummaryLayerOccurenceInformation) Name() string {
	return o.ModelName + " Layer Occurrence"
}

func (o SummaryLayerOccurenceInformation) PiePlot(title string) *charts.Pie {
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.TitleOpts{Title: title},
	)
	pie = o.PiePlotAdd(pie)
	return pie
}

type LayerAggregatedInformationSelector func(elem LayerAggregatedInformation) interface{}

func (o SummaryLayerDruationInformation) piePlotAdd(pie *charts.Pie, elemSelector LayerAggregatedInformationSelector) *charts.Pie {
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

func (o SummaryLayerDruationInformation) PiePlotAdd(pie *charts.Pie) *charts.Pie {
	return o.piePlotAdd(pie, func(elem LayerAggregatedInformation) interface{} {
		return elem.Duration
	})
}

func (o SummaryLayerOccurenceInformation) PiePlotAdd(pie *charts.Pie) *charts.Pie {
	return o.piePlotAdd(pie, func(elem LayerAggregatedInformation) interface{} {
		return elem.Occurence
	})
}

func (o SummaryLayerOccurenceInformation) WritePiePlot(path string) error {
	return writePiePlot(o, path)
}

func (o SummaryLayerOccurenceInformation) OpenPiePlot() error {
	return openPiePlot(o)
}

func (o SummaryLayerDruationInformation) WritePiePlot(path string) error {
	return writePiePlot(o, path)
}

func (o SummaryLayerDruationInformation) OpenPiePlot() error {
	return openPiePlot(o)
}

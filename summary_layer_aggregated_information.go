package evaluation

import (
	"github.com/rai-project/evaluation/writer"
	"github.com/rai-project/go-echarts/charts"
	"github.com/spf13/cast"
)

type LayerAggregatedInformation struct {
	Type       string  `json:"type,omitempty"`
	Occurences int     `json:"occurences,omitempty"`
	Duration   float64 `json:"duration,omitempty"`
}

type LayerAggregatedInformations []LayerAggregatedInformation

//easyjson:json
type SummaryLayerAggregatedInformation struct {
	SummaryBase                 `json:",inline"`
	LayerAggregatedInformations LayerAggregatedInformations `json:"layer_aggregated_informations,omitempty"`
}

func (LayerAggregatedInformation) Header(opts ...writer.Option) []string {
	return []string{
		"type",
		"occurences",
		"duration (us)",
	}
}

func (info LayerAggregatedInformation) Row(opts ...writer.Option) []string {
	return []string{
		info.Type,
		cast.ToString(info.Occurences),
		cast.ToString(info.Duration),
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
	}

	layerAggreInfos := []LayerAggregatedInformation{}
	for _, info := range exsistedLayers {
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

func (o SummaryLayerAggregatedInformation) PiePlotAdd(pie *charts.Pie) *charts.Pie {
	labels := []string{}
	data := make(map[string]interface{})
	for _, elem := range o.LayerAggregatedInformations {
		data[elem.Type] = elem.Duration
		labels = append(labels, cast.ToString(elem.Type))

	}
	pie.AddSorted("pie", data, charts.LabelTextOpts{Show: true})
	return pie
}

func (o SummaryLayerAggregatedInformation) WritePiePlot(path string) error {
	return writePiePlot(o, path)
}

func (o SummaryLayerAggregatedInformation) OpenPiePlot() error {
	return openPiePlot(o)
}

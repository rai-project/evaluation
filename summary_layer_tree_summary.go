package evaluation

import (
	"errors"
	"strings"

	"github.com/spf13/cast"
	db "upper.io/db.v3"
)

func (p Performance) LayerInformationTreeSummary(e Evaluation) (*SummaryLayerInformation, error) {
	sspans := getSpanLayersFromSpans(p.Spans())
	numSSpans := len(sspans)

	summary := &SummaryLayerInformation{
		SummaryBase:       e.summaryBase(),
		LayerInformations: LayerInformations{},
	}
	if numSSpans == 0 {
		return summary, nil
	}

	infosFullMap := make([]layerInformationMap, numSSpans)
	for ii, spans := range sspans {
		if infosFullMap[ii] == nil {
			infosFullMap[ii] = layerInformationMap{}
		}
		for _, span := range spans {
			opName := strings.ToLower(span.OperationName)
			if _, ok := infosFullMap[ii][opName]; !ok {

				infosFullMap[ii][opName] = LayerInformation{
					Name:      span.OperationName,
					Durations: []float64{},
				}
			}
			info := infosFullMap[ii][opName]
			info.Durations = append(info.Durations, cast.ToFloat64(span.Duration))
			infosFullMap[ii][opName] = info
		}
	}

	keyOrdering := []string{}
	infoMap := layerInformationMap{}
	for _, span := range sspans[0] {
		opName := strings.ToLower(span.OperationName)
		if _, ok := infoMap[opName]; !ok {
			keyOrdering = append(keyOrdering, opName)
			infoMap[opName] = LayerInformation{
				Name:      span.OperationName,
				Durations: []float64{},
			}
		}
		info := infoMap[opName]
		allDurations := [][]float64{}
		for ii := range sspans {
			allDurations = append(allDurations, infosFullMap[ii][opName].Durations)
		}
		transposedDurations := transpose(allDurations)
		durations := []float64{}
		for _, tr := range transposedDurations {
			ts := []float64{}
			for _, t := range tr {
				if t != -1 {
					ts = append(ts, t)
				}
			}
			durations = append(durations, trimmedMean(ts, DefaultTrimmedMeanFraction))
		}
		info.Durations = durations
		infoMap[opName] = info
	}

	infos := []LayerInformation{}
	for _, v := range keyOrdering {
		infos = append(infos, infoMap[v])
	}

	summary.LayerInformations = infos
	return summary, nil
}

func (e Evaluation) LayerInformationTreeSummary(perfCol *PerformanceCollection) (*SummaryLayerInformation, error) {
	perfs, err := perfCol.Find(db.Cond{"_id": e.PerformanceID})
	if err != nil {
		return nil, err
	}
	if len(perfs) != 1 {
		return nil, errors.New("expecting on performance output")
	}
	perf := perfs[0]
	return perf.LayerInformationTreeSummary(e)
}

func (es Evaluations) LayerInformationTreeSummary(perfCol *PerformanceCollection) (SummaryLayerInformations, error) {
	res := []SummaryLayerInformation{}
	for _, e := range es {
		s, err := e.LayerInformationTreeSummary(perfCol)
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

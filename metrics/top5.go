package metrics

import (
	"sort"

	"github.com/rai-project/dlframework"
)

func ClassificationTop5(features *dlframework.Features, expectedLabelIndex int) bool {
	for _, feature := range []*dlframework.Feature(*features) {
		_, ok := feature.Feature.(*dlframework.Feature_Classification)
		if !ok {
			panic("unable to convert first feature to boundingbox")
		}
	}

	sort.Sort(features)

	for _, feature := range []*dlframework.Feature(*features)[:5] {
		if int(feature.Feature.(*dlframework.Feature_Classification).Classification.Index) == expectedLabelIndex {
			return true
		}
	}
	return false
}

func Top5(features *dlframework.Features, expectedLabelIndex int) bool {
	return ClassificationTop5(features, expectedLabelIndex)
}

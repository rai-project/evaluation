package metrics

import (
	"sort"

	"github.com/rai-project/config"
	"github.com/rai-project/dlframework"
	"github.com/spf13/cast"
)

func ClassificationTop1(features *dlframework.Features, expectedLabelIndex int) bool {
	for _, feature := range []*dlframework.Feature(*features) {
		_, ok := feature.Feature.(*dlframework.Feature_Classification)
		if !ok {
			panic("unable to convert first feature to boundingbox")
		}
	}

	sort.Sort(features)

	actualLabelIndex := int([]*dlframework.Feature(*features)[0].Feature.(*dlframework.Feature_Classification).Classification.Index)
	return actualLabelIndex == expectedLabelIndex
}

func Top1(features *dlframework.Features, expectedLabelIndex int) bool {
	return ClassificationTop1(features, expectedLabelIndex)
}

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

func init() {
	config.AfterInit(func() {
		RegisterFeatureCompareFunction("Top1",
			func(actual *dlframework.Features, expected interface{}) float64 {
				expectedLabel, err := cast.ToIntE(expected)
				if err != nil {
					panic("expecting an int for second argument")
				}
				if Top1(actual, expectedLabel) {
					return 1.0
				}
				return 0.0
			})

		RegisterFeatureCompareFunction("Top5",
			func(actual *dlframework.Features, expected interface{}) float64 {
				expectedLabel, err := cast.ToIntE(expected)
				if err != nil {
					panic("expecting an int for second argument")
				}
				if Top5(actual, expectedLabel) {
					return 1.0
				}
				return 0.0
			})
	})
}

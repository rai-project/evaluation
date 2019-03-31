package metrics

import "github.com/rai-project/dlframework"

// func ClassificationTop1(features *dlframework.Features, expectedLabelIndex int) bool {
// 	for _, feature := range []*dlframework.Feature(*features) {
// 		_, ok := feature.Feature.(*dlframework.Feature_Classification)
// 		if !ok {
// 			panic("unable to convert first feature to boundingbox")
// 		}
// 	}

// 	sort.Sort(features)

// 	actualLabelIndex := int([]*dlframework.Feature(*features)[0].Feature.(*dlframework.Feature_Classification).Classification.Index)
// 	return actualLabelIndex == expectedLabelIndex
// }

func MeanAveragePrecision(features *dlframework.Features, expectedLabelIndex int) bool {
	return ClassificationTop1(features, expectedLabelIndex)
}

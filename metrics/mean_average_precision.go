package metrics

// https://github.com/alonewithyou/GoPredictor/blob/master/Measurement/metrics.go#L127
// https://github.com/ariaaan/mean-average-precision-calculation/blob/master/measure_map.py#L9
// https://forums.fast.ai/t/mean-average-precision-map/14345

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

// func MeanAveragePrecision(features *dlframework.Features, expectedLabelIndex int) bool {
// 	return ClassificationTop1(features, expectedLabelIndex)
// }

package metrics

import "math"

type metrics struct { // see http://www.cs.kau.se/pulls/hot/measurements/
	TruePositive            int // true positive
	FalsePositiveToPositive int // false-positive-to-positive
	FalseNegativeToPositive int // false-negative-to-positive
	FalseNegative           int // false negative
	TrueNegative            int // true negative
}

func (m *metrics) Add(o *metrics) {
	m.FalseNegative += o.FalseNegative
	m.FalseNegativeToPositive += o.FalseNegativeToPositive
	m.FalsePositiveToPositive += o.FalsePositiveToPositive
	m.TrueNegative += o.TrueNegative
	m.TruePositive += o.TruePositive
}

// recall = TPR = TP / (TP + FN + FPP)
func Recall(data []metrics) float64 {
	var p float64
	for i := 0; i < len(data); i++ {
		d := float64(data[i].TruePositive) / float64(data[i].TruePositive+data[i].FalseNegative+data[i].FalsePositiveToPositive)
		if !math.IsNaN(d) {
			p += d
		}
	}
	return p / float64(len(data))
}

// precision = TP / (TP + FPP + FNP)
func Precision(data []metrics) float64 {
	var p float64
	for i := 0; i < len(data); i++ {
		d := float64(data[i].TruePositive) / float64(data[i].TruePositive+data[i].FalsePositiveToPositive+data[i].FalseNegativeToPositive)
		if !math.IsNaN(d) {
			p += d
		}
	}
	return p / float64(len(data))
}

// FPR = FP / non-monitored elements = (FPP + FNP) / (TN + FNP)
func Fpr(data []metrics) float64 {
	var p float64
	for i := 0; i < len(data); i++ {
		d := float64(data[i].FalsePositiveToPositive+data[i].FalseNegativeToPositive) / float64(data[i].TrueNegative+data[i].FalseNegativeToPositive)
		if !math.IsNaN(d) {
			p += d
		}
	}
	return p / float64(len(data))
}

// F1Score = 2 * [(precision*recall) / (precision + recall)]
func F1Score(data []metrics) float64 {
	var p float64
	for i := 0; i < len(data); i++ {
		precision := float64(data[i].TruePositive) / float64(data[i].TruePositive+data[i].FalsePositiveToPositive+data[i].FalseNegativeToPositive)
		recall := float64(data[i].TruePositive) / float64(data[i].TruePositive+data[i].FalseNegative+data[i].FalsePositiveToPositive)
		if !math.IsNaN(precision) && !math.IsNaN(recall) {
			p += 2 * ((precision * recall) / (precision + recall))
		}
	}
	return p / float64(len(data))
}

// Accuracy = (TP + TN) / (everything)
func Accuracy(data []metrics) float64 {
	var p float64
	for i := 0; i < len(data); i++ {
		d := float64(data[i].TruePositive+data[i].TrueNegative) /
			float64(data[i].FalseNegative+data[i].FalseNegativeToPositive+data[i].FalsePositiveToPositive+data[i].TrueNegative+data[i].TruePositive)
		if !math.IsNaN(d) {
			p += d
		}
	}
	return p / float64(len(data))
}

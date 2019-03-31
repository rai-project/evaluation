package metrics

import "math"

/*
	Suppose pixel values in [0,1]
	refer https://en.wikipedia.org/wiki/Peak_signal-to-noise_ratio
*/
func PeakSignalToNoiseRatio(input, reference []float64) float64 {
	mse := MeanSquaredError(input, reference)
	ret := 20*math.Log10(1.0) - 10*math.Log10(float64(mse))
	return ret
}

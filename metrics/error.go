package metrics

import "math"

func MeanSquaredError(a, b []float64) float64 {
	n := len(a)
	if n != len(b) {
		panic("length not equal")
	}
	sum := 0.0
	for i := 0; i < n; i++ {
		delta := a[i] - b[i]
		sum += delta * delta
	}
	return sum / float64(n)
}

// MSPE computes the mean-square-percentage error.
func MeanSquaredPercentageError(y, yhat []float64) float64 {
	Σ := 0.0
	for i := range y {
		ε := (yhat[i] - y[i]) / y[i]
		Σ += ε * ε
	}
	return Σ / float64(len(y))
}

// NRMSE computes the normalized root-mean-square error.
//
// https://en.wikipedia.org/wiki/Root-mean-square_deviation#Normalized_root-mean-square_deviation
func NormalizedRootMeanSquaredError(y, yhat []float64) float64 {
	count := len(y)
	min, max := y[0], y[0]
	for i := 1; i < count; i++ {
		if y[i] < min {
			min = y[i]
		}
		if y[i] > max {
			max = y[i]
		}
	}
	return RootMeanSquaredError(y, yhat) / (max - min)
}

// RMSE computes the root-mean-square error.
//
// https://en.wikipedia.org/wiki/Root-mean-square_deviation
func RootMeanSquaredError(y, yhat []float64) float64 {
	return math.Sqrt(MeanSquaredError(y, yhat))
}

// RMSPE computes the root-mean-square-percentage error.
func RootMeanSquaredPercentageError(y, yhat []float64) float64 {
	return math.Sqrt(MeanSquaredPercentageError(y, yhat))
}

func RelativeAbsoluteError(a, b []float64) float64 {
	return Sum(AbsoluteError(a, b)) / Sum(AbsoluteError(a, Broadcast(Mean(a), len(a))))
}

// Compute the squared error
// This function computes the elementwise squared error for a vector
func SquaredError(a, b []float64) []float64 {
	n := len(a)
	if n != len(b) {
		panic("length not equal")
	}
	res := make([]float64, n)
	for ii := 0; ii < n; ii++ {
		diff := a[ii] - b[ii]
		res[ii] = diff * diff
	}
	return res
}

// Compute the squared log error
// This function computes the elementwise squared log error for a vector
func SquaredLogError(a, b []float64) []float64 {
	n := len(a)
	if n != len(b) {
		panic("length not equal")
	}
	res := make([]float64, n)
	for ii := 0; ii < n; ii++ {
		e := (math.Log(1+a[ii]) - math.Log(1+b[ii]))
		res[ii] = e * e
	}
	return res
}

// Compute the absolute error
// This function computes the elementwise absolute error for a vector
func AbsoluteError(a, b []float64) []float64 {
	n := len(a)
	if n != len(b) {
		panic("length not equal")
	}
	res := make([]float64, n)
	for ii := 0; ii < n; ii++ {
		res[ii] = math.Abs(a[ii] - b[ii])
	}
	return res
}

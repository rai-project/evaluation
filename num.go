package evaluation

import (
	"math"
	"sort"

	"gonum.org/v1/gonum/stat"
)

var (
	DefaultTrimmedMeanFraction = 0.2
)

func trimmedMean(data []float64, frac float64) {
	if frac == 0 {
		frac = DefaultTrimmedMeanFraction
	}

	cnt := len(data)

	sort.Float64s(data)

	start := maxInt(0, floor(cnt*frac))
	end := minInt(cnt-1, cnt-floor(cnt*frac))

	trimmed := data[start:end]

	return stat.Mean(trimmed)
}

func floor(x float64) int {
	return int(math.Floor(x))
}

func ceil(x float64) int {
	return int(math.Ceil(x))
}

func maxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func minInt(x, y int) int {
	if x > y {
		return y
	}
	return x
}

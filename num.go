package evaluation

import (
	"math"
	"sort"
)

var (
	DefaultTrimmedMeanFraction = 0.2
)

func TrimmedMean(data []float64, frac float64) float64 {

	// Sum returns the sum of the elements of the slice.
	total := func(s []float64) float64 {
		var sum float64
		for _, val := range s {
			sum += val
		}
		return sum
	}

	// Mean computes the weighted mean of the data set.
	//  sum_i {w_i * x_i} / sum_i {w_i}
	// If weights is nil then all of the weights are 1. If weights is not nil, then
	// len(x) must equal len(weights).
	mean := func(x, weights []float64) float64 {
		if weights == nil {
			return total(x) / float64(len(x))
		}
		if len(x) != len(weights) {
			panic("stat: slice length mismatch")
		}
		var (
			sumValues  float64
			sumWeights float64
		)
		for i, w := range weights {
			sumValues += w * x[i]
			sumWeights += w
		}
		return sumValues / sumWeights
	}

	if frac == 0 {
		frac = DefaultTrimmedMeanFraction
	}
	if len(data) == 0 {
		return 0
	}
	if len(data) < 3 {
		return mean(data, nil)
	}
	if len(data) == 3 {
		sort.Float64s(data)
		return data[1]
	}

	cnt := len(data)

	sort.Float64s(data)

	start := maxInt(0, floor(float64(cnt)*frac))
	end := minInt(cnt-1, cnt-floor(float64(cnt)*frac))

	// pp.Println("start = ", start, "   end = ", end)
	trimmed := data[start:end]

	ret := mean(trimmed, nil)

	return ret
}

func convertInt64SliceToFloat64Slice(in []int64) []float64 {
	ret := make([]float64, len(in))
	for i, v := range in {
		ret[i] = float64(v)
	}
	return ret
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

func uint64Transpose(a [][]uint64) [][]uint64 {
	m := len(a)
	n := len(a[0])
	aNew := make([][]uint64, n)
	for i := 0; i < n; i++ {
		aNew[i] = make([]uint64, m)
	}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			aNew[j][i] = a[i][j]
		}
	}
	return aNew
}

func transpose0(a [][]float64) [][]float64 {
	m := len(a)
	n := len(a[0])

	aNew := make([][]float64, n)
	for i := 0; i < n; i++ {
		aNew[i] = make([]float64, m)
	}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			aNew[j][i] = a[i][j]
		}
	}
	return aNew
}

func transpose(a [][]float64) [][]float64 {
	maxCols := len(a[0])
	for _, r := range a {
		maxCols = maxInt(maxCols, len(r))
	}
	r := make([][]float64, maxCols)
	for x, _ := range r {
		r[x] = make([]float64, len(a))
		for ii := range r[x] {
			r[x][ii] = -1
		}
	}
	for y, s := range a {
		for x, e := range s {
			r[x][y] = e
		}
	}
	return r
}

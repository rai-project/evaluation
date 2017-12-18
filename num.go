package evaluation

import (
	"math"
	"sort"

	"gonum.org/v1/gonum/stat"
)

var (
	DefaultTrimmedMeanFraction = 0.2
)

func trimmedMean(data []float64, frac float64) float64 {
	if frac == 0 {
		frac = DefaultTrimmedMeanFraction
	}
	if len(data) == 0 {
		return 0
	}

	cnt := len(data)

	sort.Float64s(data)

	start := maxInt(0, floor(float64(cnt)*frac))
	end := minInt(cnt-1, cnt-floor(float64(cnt)*frac))

	// pp.Println("start = ", start, "   end = ", end)
	trimmed := data[start:end]

	return stat.Mean(trimmed, nil)
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

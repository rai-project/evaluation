package evaluation

import (
	"sort"
	"time"
)

// durationMin gets the min of a slice of durations
func durationMin(input []time.Duration) time.Duration {

	if len(input) == 0 {
		return time.Duration(0)
	}
	minVal := input[0]
	for _, val := range input {
		if val < minVal {
			minVal = val
		}
	}

	return minVal
}

// durationMax gets the max of a slice of durations
func durationMax(input []time.Duration) time.Duration {

	if len(input) == 0 {
		return time.Duration(0)
	}
	maxVal := input[0]
	for _, val := range input {
		if val > maxVal {
			maxVal = val
		}
	}

	return maxVal
}

// most of what is bellow is copied from https://github.com/montanaflynn/stats
// and changed to work with time.Duration types

// durationMean gets the average of a slice of durations
func durationMean(input []time.Duration) time.Duration {

	if len(input) == 0 {
		return 0
	}

	durationSum := durationSum(input)

	return durationSum / time.Duration(len(input))
}

// durationSum adds all the durations of a slice together
func durationSum(input []time.Duration) (durationSum time.Duration) {

	if len(input) == 0 {
		return 0
	}

	// Add em up
	for _, n := range input {
		durationSum += n
	}

	return durationSum
}

// durationMedian gets the durationMedian number in a slice of durations
func durationMedian(input []time.Duration) (durationMedian time.Duration) {

	// Start by sorting a copy of the slice
	c := sortedDurationCopy(input)

	// No math is needed if there are no durations
	// For even durations we add the two middle durations
	// and divide by two using the durationMean function above
	// For odd durations we just use the middle number
	l := len(c)
	if l == 0 {
		return 0
	} else if l%2 == 0 {
		durationMedian = durationMean(c[l/2-1 : l/2+1])
	} else {
		durationMedian = c[l/2]
	}

	return durationMedian
}

// copyDurationSlice copies a slice of float64s
func copyDurationSlice(input []time.Duration) []time.Duration {
	s := make([]time.Duration, len(input))
	copy(s, input)
	return s
}

// sortedDurationCopy returns a sorted copy of float64s
func sortedDurationCopy(input []time.Duration) (copy []time.Duration) {
	copy = copyDurationSlice(input)
	sort.Slice(copy, func(ii, jj int) bool {
		return copy[ii] < copy[jj]
	})
	return
}

// durationPercentile finds the relative standing in a slice of floats
func durationPercentile(input []time.Duration, percent int64) (percentile time.Duration) {

	if len(input) == 0 {
		return 0
	}

	if percent <= 0 || percent > 100 {
		return 0
	}

	// Start by sorting a copy of the slice
	c := sortedDurationCopy(input)

	// Multiply percent by length of input
	index := (float64(percent) / 100.0) * float64(len(c))

	// Check if the index is a whole number
	if index == float64(time.Duration(index)) {

		// Convert float to int
		i := int(index)

		// Find the value at the index
		percentile = c[i-1]

	} else if index > 1 {

		// Convert float to int via truncation
		i := int(index)

		// Find the average of the index and following values
		percentile = durationMean([]time.Duration{c[i-1], c[i]})

	} else {
		return c[0]
	}

	return percentile

}

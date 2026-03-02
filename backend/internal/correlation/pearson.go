package correlation

import (
	"fmt"
	"math"
)

// PearsonCorrelation computes the Pearson r between two equal-length slices.
// Returns a value in [-1, 1]. Skips NaN values.
// Returns error if fewer than 3 valid pairs.
func PearsonCorrelation(x, y []float64) (float64, error) {
	if len(x) != len(y) {
		return 0, fmt.Errorf("slices must have equal length: %d vs %d", len(x), len(y))
	}
	// Filter out pairs where either value is NaN
	var vx, vy []float64
	for i := range x {
		if !math.IsNaN(x[i]) && !math.IsNaN(y[i]) {
			vx = append(vx, x[i])
			vy = append(vy, y[i])
		}
	}
	n := len(vx)
	if n < 3 {
		return 0, fmt.Errorf("need at least 3 valid data points, got %d", n)
	}
	mx := mean(vx)
	my := mean(vy)
	sx := stdDev(vx)
	sy := stdDev(vy)
	if sx == 0 || sy == 0 {
		return 0, fmt.Errorf("zero variance in data")
	}
	var sum float64
	for i := range vx {
		sum += (vx[i] - mx) * (vy[i] - my)
	}
	r := sum / (float64(n) * sx * sy)
	// Clamp to [-1, 1] for floating point safety
	if r > 1 {
		r = 1
	}
	if r < -1 {
		r = -1
	}
	return r, nil
}

// mean computes the arithmetic mean of a slice.
func mean(v []float64) float64 {
	if len(v) == 0 {
		return 0
	}
	var sum float64
	for _, x := range v {
		sum += x
	}
	return sum / float64(len(v))
}

// stdDev computes the population standard deviation.
func stdDev(v []float64) float64 {
	if len(v) == 0 {
		return 0
	}
	m := mean(v)
	var sumSq float64
	for _, x := range v {
		sumSq += (x - m) * (x - m)
	}
	return math.Sqrt(sumSq / float64(len(v)))
}

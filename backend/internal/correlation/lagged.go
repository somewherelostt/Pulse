package correlation

import "fmt"

// LaggedCorrelation computes the Pearson correlation between feature values
// at time t and mood values at time t+lag.
// lag=2 means: does the feature value 2 days ago predict today's mood?
// Returns the correlation, the number of valid pairs, and any error.
func LaggedCorrelation(feature []float64, mood []float64, lag int) (r float64, n int, err error) {
	if len(feature) != len(mood) {
		return 0, 0, fmt.Errorf("slice length mismatch: feature=%d mood=%d", len(feature), len(mood))
	}
	if lag < 0 || lag >= len(feature) {
		return 0, 0, fmt.Errorf("invalid lag %d for slice length %d", lag, len(feature))
	}
	total := len(feature)
	// feature[i] predicts mood[i+lag]
	var fx, my []float64
	for i := 0; i+lag < total; i++ {
		fx = append(fx, feature[i])
		my = append(my, mood[i+lag])
	}
	if len(fx) < 3 {
		return 0, len(fx), fmt.Errorf("not enough pairs after applying lag=%d: only %d", lag, len(fx))
	}
	r, err = PearsonCorrelation(fx, my)
	return r, len(fx), err
}

// BestLag searches lag values from 0 to maxLag (inclusive) and returns
// the lag with the highest absolute Pearson correlation.
func BestLag(feature []float64, mood []float64, maxLag int) (bestLag int, bestR float64, err error) {
	found := false
	for lag := 0; lag <= maxLag; lag++ {
		r, _, e := LaggedCorrelation(feature, mood, lag)
		if e != nil {
			continue
		}
		if absF(r) > absF(bestR) || !found {
			bestR = r
			bestLag = lag
			found = true
		}
	}
	if !found {
		return 0, 0, fmt.Errorf("no valid lag found in range 0-%d", maxLag)
	}
	return bestLag, bestR, nil
}

func absF(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

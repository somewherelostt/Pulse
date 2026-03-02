package correlation

import (
	"fmt"
	"math"
)

// PValue approximates the two-tailed p-value for a Pearson r with sample size n.
// Uses the t-distribution: t = r * sqrt(n-2) / sqrt(1-r^2), df = n-2.
func PValue(r float64, n int) float64 {
	if n <= 2 {
		return 1.0
	}
	if math.Abs(r) >= 1.0 {
		return 0.0
	}
	t := r * math.Sqrt(float64(n-2)) / math.Sqrt(1-r*r)
	return 2 * (1 - tCDF(math.Abs(t), float64(n-2)))
}

// IsSignificant returns true if the two-tailed p-value < alpha (typically 0.05).
func IsSignificant(r float64, n int, alpha float64) bool {
	return PValue(r, n) < alpha
}

// tCDF approximates the CDF of Student's t-distribution.
func tCDF(t, df float64) float64 {
	x := df / (df + t*t)
	return 1 - 0.5*regularizedIncompleteBeta(df/2, 0.5, x)
}

// regularizedIncompleteBeta approximates I_x(a, b).
func regularizedIncompleteBeta(a, b, x float64) float64 {
	if x <= 0 {
		return 0
	}
	if x >= 1 {
		return 1
	}
	lbeta := logBeta(a, b)
	front := math.Exp(math.Log(x)*a+math.Log(1-x)*b-lbeta) / a
	if x < (a+1)/(a+b+2) {
		return front * betaCF(a, b, x)
	}
	front2 := math.Exp(math.Log(1-x)*b+math.Log(x)*a-lbeta) / b
	return 1 - front2*betaCF(b, a, 1-x)
}

// betaCF evaluates the continued fraction representation of the incomplete beta function.
func betaCF(a, b, x float64) float64 {
	const maxIter = 200
	const eps = 3e-7
	qab := a + b
	qap := a + 1
	qam := a - 1
	c := 1.0
	d := 1 - qab*x/qap
	if math.Abs(d) < 1e-30 {
		d = 1e-30
	}
	d = 1 / d
	h := d
	for m := 1; m <= maxIter; m++ {
		m2 := 2 * float64(m)
		aa := float64(m) * (b - float64(m)) * x / ((qam + m2) * (a + m2))
		d = 1 + aa*d
		if math.Abs(d) < 1e-30 {
			d = 1e-30
		}
		c = 1 + aa/c
		if math.Abs(c) < 1e-30 {
			c = 1e-30
		}
		d = 1 / d
		h *= d * c
		aa = -(a + float64(m)) * (qab + float64(m)) * x / ((a + m2) * (qap + m2))
		d = 1 + aa*d
		if math.Abs(d) < 1e-30 {
			d = 1e-30
		}
		c = 1 + aa/c
		if math.Abs(c) < 1e-30 {
			c = 1e-30
		}
		d = 1 / d
		del := d * c
		h *= del
		if math.Abs(del-1) < eps {
			break
		}
	}
	return h
}

func logBeta(a, b float64) float64 {
	la, _ := math.Lgamma(a)
	lb, _ := math.Lgamma(b)
	lab, _ := math.Lgamma(a + b)
	return la + lb - lab
}

// PlainEnglish generates a human-readable sentence describing a correlation.
func PlainEnglish(featureName, moodMetric string, lag int, r float64) string {
	var featureDir, moodDir string
	if r > 0 {
		featureDir = "higher"
		moodDir = "higher"
	} else {
		featureDir = "higher"
		moodDir = "lower"
	}
	strength := strengthLabel(math.Abs(r))
	if lag == 0 {
		return fmt.Sprintf("%s %s correlates with %s %s on the same day (%s, r=%.2f)",
			featureName, featureDir, moodDir, moodMetric, strength, r)
	}
	return fmt.Sprintf("%s %s %d day(s) ago predicts %s %s (%s, r=%.2f)",
		featureName, featureDir, lag, moodDir, moodMetric, strength, r)
}

func strengthLabel(absR float64) string {
	switch {
	case absR >= 0.7:
		return "strong"
	case absR >= 0.4:
		return "moderate"
	default:
		return "weak"
	}
}

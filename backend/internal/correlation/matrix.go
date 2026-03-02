package correlation

import (
	"math"
	"sort"
)

// CorrelationResult holds the result of a single feature vs mood correlation.
type CorrelationResult struct {
	FeatureName  string  `json:"feature_name"`
	MoodMetric   string  `json:"mood_metric"`
	BestLag      int     `json:"best_lag"`
	Correlation  float64 `json:"correlation"`
	PValue       float64 `json:"p_value"`
	SampleSize   int     `json:"sample_size"`
	Significant  bool    `json:"significant"`
	Direction    string  `json:"direction"`
	PlainEnglish string  `json:"plain_english"`
}

// FeatureNames is the ordered list of calendar features to correlate.
var FeatureNames = []string{
	"meeting_density_pct",
	"fragmentation_score",
	"avg_focus_block_mins",
	"after_hours_mins",
	"back_to_back_count",
	"avg_recovery_mins",
	"attendee_avg",
	"solo_time_pct",
	// Layer 2 sleep features (included when available)
	"total_sleep_hours",
	"sleep_efficiency_pct",
	"hrv_rmssd",
	"circadian_health_score",
	"sleep_debt_rolling",
	"bedtime_variance_mins",
	"social_jetlag_mins",
}

// MoodMetrics is the list of mood dimensions to target.
var MoodMetrics = []string{"mood_score", "energy_score"}

// BuildMatrix computes the correlation for every feature against every mood metric
// across lags 0-7, and returns results sorted by absolute correlation descending.
// features and mood are maps of name -> []float64 (oldest first).
func BuildMatrix(features map[string][]float64, mood map[string][]float64) ([]CorrelationResult, error) {
	var results []CorrelationResult
	for _, fn := range FeatureNames {
		fv, ok := features[fn]
		if !ok {
			continue
		}
		for _, mn := range MoodMetrics {
			mv, ok := mood[mn]
			if !ok {
				continue
			}
			bl, br, err := BestLag(fv, mv, 7)
			if err != nil {
				continue
			}
			_, n, _ := LaggedCorrelation(fv, mv, bl)
			pval := PValue(br, n)
			sig := IsSignificant(br, n, 0.05)
			dir := "positive"
			if br < 0 {
				dir = "negative"
			}
			plain := PlainEnglish(fn, mn, bl, br)
			results = append(results, CorrelationResult{
				FeatureName:  fn,
				MoodMetric:   mn,
				BestLag:      bl,
				Correlation:  br,
				PValue:       pval,
				SampleSize:   n,
				Significant:  sig,
				Direction:    dir,
				PlainEnglish: plain,
			})
		}
	}
	sort.Slice(results, func(i, j int) bool {
		return math.Abs(results[i].Correlation) > math.Abs(results[j].Correlation)
	})
	return results, nil
}

package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// DaySnapshot holds one day's worth of behavioral and mood data for LLM analysis.
type DaySnapshot struct {
	Date               string   `json:"date"`
	MoodScore          *float64 `json:"mood_score,omitempty"`
	EnergyScore        *float64 `json:"energy_score,omitempty"`
	MeetingDensityPct  *float64 `json:"meeting_density_pct,omitempty"`
	FragmentationScore *float64 `json:"fragmentation_score,omitempty"`
	AfterHoursMins     *float64 `json:"after_hours_mins,omitempty"`
	BackToBackCount    *int     `json:"back_to_back_count,omitempty"`
	AvgFocusBlockMins  *float64 `json:"avg_focus_block_mins,omitempty"`
	AvgRecoveryMins    *float64 `json:"avg_recovery_mins,omitempty"`
	AttendeeAvg        *float64 `json:"attendee_avg,omitempty"`
}

// PatternAnalysisInput is the input to AnalyzePatterns.
type PatternAnalysisInput struct {
	UserID   string
	Days     []DaySnapshot
	Timezone string
}

// DetectedPattern is a single pattern identified by the LLM.
type DetectedPattern struct {
	PatternID        string   `json:"pattern_id"`
	FeaturesInvolved []string `json:"features_involved"`
	Direction        string   `json:"direction"`
	LagDays          int      `json:"lag_days"`
	Confidence       float64  `json:"confidence"`
	PlainEnglish     string   `json:"plain_english"`
	Severity         string   `json:"severity"`
}

// PatternAnalysisOutput is the structured result from the LLM.
type PatternAnalysisOutput struct {
	Patterns        []DetectedPattern `json:"patterns"`
	Summary         string            `json:"summary"`
	Recommendation  string            `json:"recommendation"`
	DataQualityNote string            `json:"data_quality_note"`
	GeneratedAt     time.Time         `json:"generated_at"`
	ModelUsed       string            `json:"model_used"`
}

// AnalyzePatterns sends behavioral data to the LLM and returns detected patterns.
func AnalyzePatterns(ctx context.Context, client LLMClient, input PatternAnalysisInput) (*PatternAnalysisOutput, error) {
	// Load system prompt from file; fall back to inline if not found
	systemPromptBytes, err := os.ReadFile("internal/llm/prompts/pattern_analysis.txt")
	if err != nil {
		systemPromptBytes = []byte(`You are a behavioral health analyst. Analyze the data and return ONLY valid JSON with a "patterns" array, "summary", "recommendation", and "data_quality_note" fields.`)
	}

	userPrompt := buildAnalysisPrompt(input)

	req := CompletionRequest{
		SystemPrompt: string(systemPromptBytes),
		UserPrompt:   userPrompt,
		MaxTokens:    2000,
		Temperature:  0.3,
		JSONMode:     true,
	}

	resp, err := client.Complete(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("LLM call failed: %w", err)
	}

	var out PatternAnalysisOutput
	if jsonErr := json.Unmarshal([]byte(resp.Content), &out); jsonErr != nil {
		// Retry once with explicit JSON reminder
		req.UserPrompt += "\n\nCRITICAL: Return ONLY valid JSON. No markdown code blocks, no prose outside JSON."
		resp, err = client.Complete(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("LLM retry failed: %w", err)
		}
		if jsonErr2 := json.Unmarshal([]byte(resp.Content), &out); jsonErr2 != nil {
			return nil, fmt.Errorf("failed to parse LLM response as JSON: %w", jsonErr2)
		}
	}

	out.GeneratedAt = time.Now()
	out.ModelUsed = resp.Model
	return &out, nil
}

func buildAnalysisPrompt(input PatternAnalysisInput) string {
	var sb strings.Builder
	sb.WriteString("Analyze these behavioral data points and identify patterns predicting mood changes.\n\n")
	sb.WriteString(fmt.Sprintf("User timezone: %s\n\n", input.Timezone))
	sb.WriteString("Date       | Mood | Energy | Density%% | Fragm | AfterHrs | B2B | FocusBlk | Recovery | Attendees\n")
	sb.WriteString("-----------|------|--------|----------|-------|----------|-----|----------|----------|-----------\n")
	for _, d := range input.Days {
		sb.WriteString(fmt.Sprintf("%-10s | %-4s | %-6s | %-8s | %-5s | %-8s | %-3s | %-8s | %-8s | %s\n",
			d.Date,
			fmtF1(d.MoodScore),
			fmtF1(d.EnergyScore),
			fmtPct(d.MeetingDensityPct),
			fmtF2(d.FragmentationScore),
			fmtF1(d.AfterHoursMins),
			fmtIntP(d.BackToBackCount),
			fmtF1(d.AvgFocusBlockMins),
			fmtF1(d.AvgRecoveryMins),
			fmtF1(d.AttendeeAvg),
		))
	}
	return sb.String()
}

func fmtF1(v *float64) string {
	if v == nil {
		return "-"
	}
	return fmt.Sprintf("%.1f", *v)
}

func fmtF2(v *float64) string {
	if v == nil {
		return "-"
	}
	return fmt.Sprintf("%.2f", *v)
}

func fmtPct(v *float64) string {
	if v == nil {
		return "-"
	}
	return fmt.Sprintf("%.0f%%", *v*100)
}

func fmtIntP(v *int) string {
	if v == nil {
		return "-"
	}
	return fmt.Sprintf("%d", *v)
}

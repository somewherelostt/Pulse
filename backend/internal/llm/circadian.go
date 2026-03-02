package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// CircadianInput holds the data for one day passed to the circadian LLM.
type CircadianInput struct {
	Date                 string   `json:"date"`
	SleepDurationMins    *float64 `json:"sleep_duration_mins,omitempty"`
	SleepEfficiencyPct   *float64 `json:"sleep_efficiency_pct,omitempty"`
	SleepDebtMins        *float64 `json:"sleep_debt_mins,omitempty"`
	MidSleepHour         *float64 `json:"mid_sleep_hour,omitempty"`
	RhythmConsistencyPct *float64 `json:"rhythm_consistency_pct,omitempty"`
	SocialJetlagMins     *float64 `json:"social_jetlag_mins,omitempty"`
	REMPct               *float64 `json:"rem_pct,omitempty"`
	DeepPct              *float64 `json:"deep_pct,omitempty"`
	HRV                  *float64 `json:"hrv_rmssd,omitempty"`
	RestingHR            *float64 `json:"resting_hr,omitempty"`
	SleepScore           *float64 `json:"sleep_score,omitempty"`
	LightHygieneScore    *float64 `json:"light_hygiene_score,omitempty"`
}

// CircadianNarrativeOutput is the JSON the LLM returns.
type CircadianNarrativeOutput struct {
	Narrative        string                  `json:"narrative"`
	Interventions    []CircadianIntervention `json:"interventions"`
	DataQualityNote  string                  `json:"data_quality_note,omitempty"`
	GeneratedAt      time.Time               `json:"generated_at"`
	ModelUsed        string                  `json:"model_used"`
}

// CircadianIntervention is one actionable item.
type CircadianIntervention struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"` // "high"|"medium"|"low"
}

// GenerateCircadianNarrative calls the LLM to produce a weekly sleep health narrative.
func GenerateCircadianNarrative(ctx context.Context, client LLMClient, days []CircadianInput) (*CircadianNarrativeOutput, error) {
	systemPrompt, err := loadCircadianPrompt()
	if err != nil {
		return nil, fmt.Errorf("load prompt: %w", err)
	}

	dataJSON, err := json.MarshalIndent(days, "", "  ")
	if err != nil {
		return nil, err
	}
	userPrompt := fmt.Sprintf("Here is the sleep data for the past %d days:\n\n%s\n\nGenerate the JSON narrative.", len(days), dataJSON)

	resp, err := client.Complete(ctx, CompletionRequest{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
		MaxTokens:    1200,
		Temperature:  0.4,
		JSONMode:     true,
	})
	if err != nil {
		return nil, fmt.Errorf("llm complete: %w", err)
	}

	var out CircadianNarrativeOutput
	if err := json.Unmarshal([]byte(resp.Content), &out); err != nil {
		// Retry with more explicit instruction
		userPrompt2 := userPrompt + "\n\nIMPORTANT: Respond with ONLY valid JSON, no other text."
		resp2, err2 := client.Complete(ctx, CompletionRequest{
			SystemPrompt: systemPrompt,
			UserPrompt:   userPrompt2,
			MaxTokens:    1200,
			Temperature:  0.2,
			JSONMode:     true,
		})
		if err2 != nil {
			return nil, fmt.Errorf("llm retry: %w", err2)
		}
		if err3 := json.Unmarshal([]byte(resp2.Content), &out); err3 != nil {
			return nil, fmt.Errorf("parse json: %w", err3)
		}
		resp = resp2
	}

	out.GeneratedAt = time.Now().UTC()
	out.ModelUsed = resp.Model
	return &out, nil
}

func loadCircadianPrompt() (string, error) {
	// Try relative paths
	candidates := []string{
		"internal/llm/prompts/circadian_narrative.txt",
		"../llm/prompts/circadian_narrative.txt",
	}
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		dir := filepath.Dir(filename)
		candidates = append([]string{filepath.Join(dir, "prompts/circadian_narrative.txt")}, candidates...)
	}
	for _, p := range candidates {
		b, err := os.ReadFile(p)
		if err == nil {
			return strings.TrimSpace(string(b)), nil
		}
	}
	// Embedded fallback
	return circadianPromptFallback, nil
}

const circadianPromptFallback = `You are a circadian health analyst. Analyze the user's sleep data and return JSON with keys: narrative (string), interventions (array of {title,description,priority}), data_quality_note (string). Be specific, actionable, and evidence-informed.`

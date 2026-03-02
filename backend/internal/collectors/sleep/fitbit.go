package sleep

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const fitbitBaseURL = "https://api.fitbit.com/1.2/user/-"

// FitbitCollector fetches sleep data using a Fitbit OAuth2 access token.
type FitbitCollector struct {
	AccessToken string
	Client      *http.Client
}

func NewFitbitCollector(accessToken string) *FitbitCollector {
	return &FitbitCollector{
		AccessToken: accessToken,
		Client:      &http.Client{Timeout: 20 * time.Second},
	}
}

type fitbitSleepLog struct {
	Sleep []struct {
		LogID         int64  `json:"logId"`
		StartTime     string `json:"startTime"`
		EndTime       string `json:"endTime"`
		DateOfSleep   string `json:"dateOfSleep"`
		Duration      int    `json:"duration"` // ms
		Efficiency    int    `json:"efficiency"`
		MinutesAsleep int    `json:"minutesAsleep"`
		MinutesAwake  int    `json:"minutesAwake"`
		Levels        struct {
			Summary struct {
				Rem   struct{ Minutes int `json:"minutes"` } `json:"rem"`
				Deep  struct{ Minutes int `json:"minutes"` } `json:"deep"`
				Light struct{ Minutes int `json:"minutes"` } `json:"light"`
				Wake  struct{ Minutes int `json:"minutes"` } `json:"wake"`
			} `json:"summary"`
		} `json:"levels"`
	} `json:"sleep"`
}

func (f *FitbitCollector) FetchSessions(from, to time.Time) ([]SleepSession, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	var sessions []SleepSession
	cur := from
	for !cur.After(to) {
		batch, err := f.fetchDay(ctx, cur)
		if err != nil {
			// Non-fatal: log and continue
			cur = cur.AddDate(0, 0, 1)
			continue
		}
		sessions = append(sessions, batch...)
		cur = cur.AddDate(0, 0, 1)
	}
	return sessions, nil
}

func (f *FitbitCollector) fetchDay(ctx context.Context, date time.Time) ([]SleepSession, error) {
	url := fmt.Sprintf("%s/sleep/date/%s.json", fitbitBaseURL, date.Format("2006-01-02"))
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+f.AccessToken)
	req.Header.Set("Accept-Language", "en_US")

	resp, err := f.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 429 {
		return nil, fmt.Errorf("fitbit rate limited")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("fitbit API %d", resp.StatusCode)
	}

	var doc fitbitSleepLog
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return nil, err
	}

	sessions := make([]SleepSession, 0, len(doc.Sleep))
	for _, sl := range doc.Sleep {
		d, err := time.Parse("2006-01-02", sl.DateOfSleep)
		if err != nil {
			continue
		}
		s := SleepSession{
			Provider:    "fitbit",
			SessionDate: d,
			SourceID:    fmt.Sprintf("%d", sl.LogID),
		}
		if bt, err := time.Parse("2006-01-02T15:04:05.000", sl.StartTime); err == nil {
			s.Bedtime = &bt
		}
		if wt, err := time.Parse("2006-01-02T15:04:05.000", sl.EndTime); err == nil {
			s.WakeTime = &wt
		}
		if sl.MinutesAsleep > 0 {
			s.TotalMins = &sl.MinutesAsleep
		}
		if sl.MinutesAwake > 0 {
			s.AwakeMins = &sl.MinutesAwake
		}
		rem := sl.Levels.Summary.Rem.Minutes
		if rem > 0 {
			s.REMMins = &rem
		}
		deep := sl.Levels.Summary.Deep.Minutes
		if deep > 0 {
			s.DeepMins = &deep
		}
		light := sl.Levels.Summary.Light.Minutes
		if light > 0 {
			s.LightMins = &light
		}
		if sl.Efficiency > 0 {
			s.SleepScore = &sl.Efficiency
		}
		sessions = append(sessions, s)
	}
	return sessions, nil
}

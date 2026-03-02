package sleep

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const ouraBaseURL = "https://api.ouraring.com/v2"

// OuraCollector fetches sleep data using a personal access token.
type OuraCollector struct {
	Token  string
	Client *http.Client
}

func NewOuraCollector(token string) *OuraCollector {
	return &OuraCollector{
		Token:  token,
		Client: &http.Client{Timeout: 20 * time.Second},
	}
}

type ouraSleepDoc struct {
	Data []struct {
		ID            string  `json:"id"`
		Day           string  `json:"day"` // "2024-01-15"
		Bedtime       string  `json:"bedtime_start"`
		WakeTime      string  `json:"bedtime_end"`
		TotalSleep    int     `json:"total_sleep_duration"` // seconds
		REMSleep      int     `json:"rem_sleep_duration"`
		DeepSleep     int     `json:"deep_sleep_duration"`
		LightSleep    int     `json:"light_sleep_duration"`
		Awake         int     `json:"awake_duration"`
		Score         int     `json:"score"`
		AverageHRV    float64 `json:"average_hrv"`
		LowestHeartRate int  `json:"lowest_heart_rate"`
	} `json:"data"`
}

func (o *OuraCollector) FetchSessions(from, to time.Time) ([]SleepSession, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	url := fmt.Sprintf("%s/usercollection/sleep?start_date=%s&end_date=%s",
		ouraBaseURL,
		from.Format("2006-01-02"),
		to.Format("2006-01-02"),
	)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+o.Token)

	resp, err := o.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("oura API %d", resp.StatusCode)
	}

	var doc ouraSleepDoc
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return nil, err
	}

	sessions := make([]SleepSession, 0, len(doc.Data))
	for _, d := range doc.Data {
		date, err := time.Parse("2006-01-02", d.Day)
		if err != nil {
			continue
		}
		s := SleepSession{
			Provider:    "oura",
			SessionDate: date,
			SourceID:    d.ID,
		}
		if bt, err := time.Parse(time.RFC3339, d.Bedtime); err == nil {
			s.Bedtime = &bt
		}
		if wt, err := time.Parse(time.RFC3339, d.WakeTime); err == nil {
			s.WakeTime = &wt
		}
		if d.TotalSleep > 0 {
			v := d.TotalSleep / 60
			s.TotalMins = &v
		}
		if d.REMSleep > 0 {
			v := d.REMSleep / 60
			s.REMMins = &v
		}
		if d.DeepSleep > 0 {
			v := d.DeepSleep / 60
			s.DeepMins = &v
		}
		if d.LightSleep > 0 {
			v := d.LightSleep / 60
			s.LightMins = &v
		}
		if d.Awake > 0 {
			v := d.Awake / 60
			s.AwakeMins = &v
		}
		if d.Score > 0 {
			s.SleepScore = &d.Score
		}
		if d.AverageHRV > 0 {
			s.HRV = &d.AverageHRV
		}
		if d.LowestHeartRate > 0 {
			hr := float64(d.LowestHeartRate)
			s.RestingHR = &hr
		}
		sessions = append(sessions, s)
	}
	return sessions, nil
}

package sleep

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/oauth2"
)

const googleFitBaseURL = "https://www.googleapis.com/fitness/v1/users/me"

// GoogleFitCollector fetches sleep data from Google Fit REST API.
type GoogleFitCollector struct {
	Token  *oauth2.Token
	Config *oauth2.Config
	Client *http.Client
}

func NewGoogleFitCollector(token *oauth2.Token, cfg *oauth2.Config) *GoogleFitCollector {
	ctx := context.Background()
	return &GoogleFitCollector{
		Token:  token,
		Config: cfg,
		Client: cfg.Client(ctx, token),
	}
}

type fitAggregateBucket struct {
	StartTimeMillis string `json:"startTimeMillis"`
	EndTimeMillis   string `json:"endTimeMillis"`
	Dataset         []struct {
		Point []struct {
			StartTimeNanos string `json:"startTimeNanos"`
			EndTimeNanos   string `json:"endTimeNanos"`
			Value          []struct {
				IntVal    int     `json:"intVal"`
				FpVal     float64 `json:"fpVal"`
				StringVal string  `json:"stringVal"`
			} `json:"value"`
		} `json:"point"`
	} `json:"dataset"`
}

type fitAggregateResp struct {
	Bucket []fitAggregateBucket `json:"bucket"`
}

func (g *GoogleFitCollector) FetchSessions(from, to time.Time) ([]SleepSession, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	startMs := from.UnixMilli()
	endMs := to.UnixMilli()

	body := fmt.Sprintf(`{
		"aggregateBy": [{"dataTypeName": "com.google.sleep.segment"}],
		"bucketByTime": {"durationMillis": 86400000},
		"startTimeMillis": %d,
		"endTimeMillis": %d
	}`, startMs, endMs)

	req, err := http.NewRequestWithContext(ctx, "POST",
		googleFitBaseURL+"/dataset:aggregate",
		stringReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("google fit API %d", resp.StatusCode)
	}

	var doc fitAggregateResp
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return nil, err
	}

	var sessions []SleepSession
	for _, bucket := range doc.Bucket {
		startMS, _ := strconv.ParseInt(bucket.StartTimeMillis, 10, 64)
		date := time.UnixMilli(startMS).UTC().Truncate(24 * time.Hour)

		var totalMins, lightMins, deepMins, remMins int
		var bedtime, wakeTime *time.Time
		for _, ds := range bucket.Dataset {
			for _, pt := range ds.Point {
				startNs, _ := strconv.ParseInt(pt.StartTimeNanos, 10, 64)
				endNs, _ := strconv.ParseInt(pt.EndTimeNanos, 10, 64)
				startT := time.Unix(0, startNs).UTC()
				endT := time.Unix(0, endNs).UTC()
				durationMins := int(endT.Sub(startT).Minutes())
				sleepType := 0
				if len(pt.Value) > 0 {
					sleepType = pt.Value[0].IntVal
				}
				// Google Fit sleep types: 1=awake,2=sleep,3=OOS,4=light,5=deep,6=REM
				switch sleepType {
				case 4:
					lightMins += durationMins
					totalMins += durationMins
				case 5:
					deepMins += durationMins
					totalMins += durationMins
				case 6:
					remMins += durationMins
					totalMins += durationMins
				case 2:
					totalMins += durationMins
				}
				if bedtime == nil || startT.Before(*bedtime) {
					t := startT
					bedtime = &t
				}
				if wakeTime == nil || endT.After(*wakeTime) {
					t := endT
					wakeTime = &t
				}
			}
		}
		if totalMins == 0 {
			continue
		}
		s := SleepSession{
			Provider:    "googlefit",
			SessionDate: date,
			Bedtime:     bedtime,
			WakeTime:    wakeTime,
		}
		s.TotalMins = &totalMins
		if lightMins > 0 {
			s.LightMins = &lightMins
		}
		if deepMins > 0 {
			s.DeepMins = &deepMins
		}
		if remMins > 0 {
			s.REMMins = &remMins
		}
		sessions = append(sessions, s)
	}
	return sessions, nil
}

type stringReaderT struct {
	s   string
	pos int
}

func (r *stringReaderT) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.s) {
		return 0, fmt.Errorf("EOF")
	}
	n = copy(p, r.s[r.pos:])
	r.pos += n
	return
}

func stringReader(s string) *stringReaderT {
	return &stringReaderT{s: s}
}

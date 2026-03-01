package google

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/api/calendar/v3"
)

func FetchEvents(ctx context.Context, svc *calendar.Service, lookbackDays int) ([]Event, error) {
	timeMin := time.Now().AddDate(0, 0, -lookbackDays).Format(time.RFC3339)
	timeMax := time.Now().AddDate(0, 0, 1).Format(time.RFC3339)

	var all []Event
	pageToken := ""
	for {
		call := svc.Events.List("primary").
			TimeMin(timeMin).
			TimeMax(timeMax).
			SingleEvents(true).
			OrderBy("startTime").
			MaxResults(2500)
		if pageToken != "" {
			call = call.PageToken(pageToken)
		}
		list, err := call.Do()
		if err != nil {
			return nil, fmt.Errorf("list events: %w", err)
		}
		for _, e := range list.Items {
			ev, err := toEvent(e)
			if err != nil {
				slog.Debug("skip event", "id", e.Id, "err", err)
				continue
			}
			all = append(all, ev)
		}
		pageToken = list.NextPageToken
		if pageToken == "" {
			break
		}
	}
	return all, nil
}

func toEvent(e *calendar.Event) (Event, error) {
	var start, end time.Time
	if e.Start.DateTime != "" {
		var err error
		start, err = time.Parse(time.RFC3339, e.Start.DateTime)
		if err != nil {
			return Event{}, err
		}
		end, err = time.Parse(time.RFC3339, e.End.DateTime)
		if err != nil {
			return Event{}, err
		}
	} else {
		start, _ = time.Parse("2006-01-02", e.Start.Date)
		end, _ = time.Parse("2006-01-02", e.End.Date)
	}
	title := e.Summary
	if title == "" {
		title = "(No title)"
	}
	attendees := 0
	if e.Attendees != nil {
		attendees = len(e.Attendees)
	}
	recurring := e.RecurringEventId != ""
	return Event{
		ID:          e.Id,
		Title:       title,
		Start:       start,
		End:         end,
		Attendees:   attendees,
		IsAllDay:    e.Start.DateTime == "",
		IsRecurring: recurring,
	}, nil
}

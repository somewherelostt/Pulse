package google

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func Config(clientID, clientSecret, redirectURI string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		Scopes:       []string{calendar.CalendarReadonlyScope},
		Endpoint:     oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}
}

func AuthCodeURL(cfg *oauth2.Config, state string) string {
	return cfg.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce, oauth2.Prompt("consent"))
}

func Exchange(ctx context.Context, cfg *oauth2.Config, code string) (*oauth2.Token, error) {
	return cfg.Exchange(ctx, code)
}

func NewCalendarService(ctx context.Context, token *oauth2.Token, cfg *oauth2.Config) (*calendar.Service, error) {
	client := cfg.Client(ctx, token)
	svc, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("calendar service: %w", err)
	}
	return svc, nil
}

func TokenSource(ctx context.Context, cfg *oauth2.Config, token *oauth2.Token) oauth2.TokenSource {
	return cfg.TokenSource(ctx, token)
}

// RefreshToken refreshes the token if needed; call before calendar API calls.
func RefreshToken(ctx context.Context, cfg *oauth2.Config, token *oauth2.Token) (*oauth2.Token, error) {
	ts := cfg.TokenSource(ctx, token)
	newTok, err := ts.Token()
	if err != nil {
		return nil, err
	}
	return newTok, nil
}

// HTTPClient returns an HTTP client that uses the given token (with refresh).
func HTTPClient(ctx context.Context, cfg *oauth2.Config, token *oauth2.Token) *http.Client {
	return cfg.Client(ctx, token)
}

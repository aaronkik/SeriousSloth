package twitch

import (
	"context"
	"emotes-service/src/adapters/secondary/parameter"
	"emotes-service/src/environment"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/sync/errgroup"
)

func GetAccessToken(ctx context.Context) (string, error) {
	var twitchClientId, twitchClientSecret string

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() (err error) {
		twitchClientId, err = parameter.GetSecret(gCtx, environment.GetOrFatal("TWITCH_CLIENT_ID_PARAM_ARN"))
		if err != nil {
			slog.ErrorContext(gCtx, "Error getting client id", "error", err)
		}
		return err
	})

	g.Go(func() (err error) {
		twitchClientSecret, err = parameter.GetSecret(gCtx, environment.GetOrFatal("TWITCH_CLIENT_SECRET_PARAM_ARN"))
		if err != nil {
			slog.ErrorContext(gCtx, "Error getting client secret", "error", err)
		}
		return err
	})

	if err := g.Wait(); err != nil {
		return "", err
	}

	oauth2Config := &clientcredentials.Config{
		ClientID:     twitchClientId,
		ClientSecret: twitchClientSecret,
		TokenURL:     environment.GetOrFatal("TWITCH_OAUTH_ENDPOINT"),
	}

	tokenCtx := context.WithValue(ctx, oauth2.HTTPClient, &http.Client{
		Timeout:   time.Second * 10,
		Transport: newrelic.NewRoundTripper(nil),
	})
	token, err := oauth2Config.Token(tokenCtx)
	if err != nil {
		slog.Error("Error getting access token", "error", err)
		return "", err
	}

	return token.AccessToken, nil
}

type GlobalEmote struct {
	Format []string `json:"format"`
	ID     string   `json:"id"`
	Images struct {
		URL1X string `json:"url_1x"`
		URL2X string `json:"url_2x"`
		URL4X string `json:"url_4x"`
	} `json:"images"`
	Name      string   `json:"name"`
	Scale     []string `json:"scale"`
	ThemeMode []string `json:"theme_mode"`
}

type GlobalEmotesResponse struct {
	Data     []GlobalEmote `json:"data"`
	Template string        `json:"template"`
}

func GetGlobalEmotes(ctx context.Context, accessToken string) ([]GlobalEmote, error) {
	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: newrelic.NewRoundTripper(nil),
	}

	req, err := http.NewRequestWithContext(ctx, "GET", environment.GetOrFatal("TWITCH_GLOBAL_EMOTES_ENDPOINT"), nil)
	if err != nil {
		slog.Error("Error creating global emotes request", "error", err)
		return nil, err
	}

	twitchClientId, err := parameter.GetSecret(ctx, environment.GetOrFatal("TWITCH_CLIENT_ID_PARAM_ARN"))
	if err != nil {
		slog.Error("Error getting client id", "error", err)
		return nil, err
	}

	req.Header.Set("Client-ID", twitchClientId)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Error getting global emotes", "error", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Unexpected status code", "status", resp.StatusCode, "response", resp)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	globalEmotesResponse := GlobalEmotesResponse{}
	err = json.NewDecoder(resp.Body).Decode(&globalEmotesResponse)
	if err != nil {
		slog.Error("Error decoding global emotes response", "error", err)
		return nil, err
	}

	slog.Info("Got global emotes", "body", globalEmotesResponse)
	return globalEmotesResponse.Data, nil
}

func GetChannelEmotes(ctx context.Context, accessToken, broadcasterId string) ([]GlobalEmote, error) {
	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: newrelic.NewRoundTripper(nil),
	}

	base := environment.GetOrFatal("TWITCH_CHANNEL_EMOTES_ENDPOINT")
	requestUrl := fmt.Sprintf("%s?broadcaster_id=%s", base, url.QueryEscape(broadcasterId))

	req, err := http.NewRequestWithContext(ctx, "GET", requestUrl, nil)
	if err != nil {
		slog.ErrorContext(ctx, "Error creating channel emotes request", "error", err)
		return nil, err
	}

	twitchClientId, err := parameter.GetSecret(ctx, environment.GetOrFatal("TWITCH_CLIENT_ID_PARAM_ARN"))
	if err != nil {
		slog.ErrorContext(ctx, "Error getting client id", "error", err)
		return nil, err
	}

	req.Header.Set("Client-ID", twitchClientId)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := client.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting channel emotes", "broadcasterId", broadcasterId, "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.ErrorContext(ctx, "Unexpected status code from channel emotes endpoint", "status", resp.StatusCode, "broadcasterId", broadcasterId)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var emotesResponse GlobalEmotesResponse
	if err := json.NewDecoder(resp.Body).Decode(&emotesResponse); err != nil {
		slog.ErrorContext(ctx, "Error decoding channel emotes response", "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "Got channel emotes", "broadcasterId", broadcasterId, "count", len(emotesResponse.Data))
	return emotesResponse.Data, nil
}

type TwitchUser struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	ProfileImageURL string `json:"profile_image_url"`
}

type GetUsersResponse struct {
	Data []TwitchUser `json:"data"`
}

// GetUserById fetches a single Twitch user by id. Returns (nil, nil) when
// Twitch returns 200 with an empty data array (i.e. user not found).
func GetUserById(ctx context.Context, accessToken, twitchId string) (*TwitchUser, error) {
	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: newrelic.NewRoundTripper(nil),
	}

	base := environment.GetOrFatal("TWITCH_USERS_ENDPOINT")
	requestUrl := fmt.Sprintf("%s?id=%s", base, url.QueryEscape(twitchId))

	req, err := http.NewRequestWithContext(ctx, "GET", requestUrl, nil)
	if err != nil {
		slog.ErrorContext(ctx, "Error creating users request", "error", err)
		return nil, err
	}

	twitchClientId, err := parameter.GetSecret(ctx, environment.GetOrFatal("TWITCH_CLIENT_ID_PARAM_ARN"))
	if err != nil {
		slog.ErrorContext(ctx, "Error getting client id", "error", err)
		return nil, err
	}

	req.Header.Set("Client-ID", twitchClientId)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := client.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting twitch user", "twitchId", twitchId, "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.ErrorContext(ctx, "Unexpected status code from twitch users endpoint", "status", resp.StatusCode, "twitchId", twitchId)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var usersResponse GetUsersResponse
	if err := json.NewDecoder(resp.Body).Decode(&usersResponse); err != nil {
		slog.ErrorContext(ctx, "Error decoding twitch users response", "error", err)
		return nil, err
	}

	if len(usersResponse.Data) == 0 {
		return nil, nil
	}

	return &usersResponse.Data[0], nil
}

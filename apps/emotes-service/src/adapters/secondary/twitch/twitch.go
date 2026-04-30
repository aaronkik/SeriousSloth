package twitch

import (
	"context"
	"emotes-service/src/adapters/secondary/parameter"
	"emotes-service/src/environment"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/oauth2/clientcredentials"
)

func GetAccessToken(ctx context.Context) (string, error) {
	twitchClientId, err := parameter.GetSecret(ctx, environment.GetOrFatal("TWITCH_CLIENT_ID_PARAM_ARN"))
	if err != nil {
		slog.Error("Error getting client id", "error", err)
		return "", err
	}

	twitchClientSecret, err := parameter.GetSecret(ctx, environment.GetOrFatal("TWITCH_CLIENT_SECRET_PARAM_ARN"))
	if err != nil {
		slog.Error("Error getting client secret", "error", err)
		return "", err
	}

	oauth2Config := &clientcredentials.Config{
		ClientID:     twitchClientId,
		ClientSecret: twitchClientSecret,
		TokenURL:     environment.GetOrFatal("TWITCH_OAUTH_ENDPOINT"),
	}

	token, err := oauth2Config.Token(context.Background())
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
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", environment.GetOrFatal("TWITCH_GLOBAL_EMOTES_ENDPOINT"), nil)
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

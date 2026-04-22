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

func GetAccessToken() (string, error) {
	twitchClientId, err := parameter.GetSecret(environment.GetOrFatal("TWITCH_CLIENT_ID_PARAM_ARN"))
	if err != nil {
		slog.Error("Error getting client id", "error", err)
		return "", err
	}

	twitchClientSecret, err := parameter.GetSecret(environment.GetOrFatal("TWITCH_CLIENT_SECRET_PARAM_ARN"))
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

func GetGlobalEmotes(accessToken string) (GlobalEmotesResponse, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/chat/emotes/global", nil)
	if err != nil {
		slog.Error("Error creating global emotes request", "error", err)
		return GlobalEmotesResponse{}, err
	}

	twitchClientId, err := parameter.GetSecret(environment.GetOrFatal("TWITCH_CLIENT_ID_PARAM_ARN"))
	if err != nil {
		slog.Error("Error getting client id", "error", err)
		return GlobalEmotesResponse{}, err
	}

	req.Header.Set("Client-ID", twitchClientId)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Error getting global emotes", "error", err)
		return GlobalEmotesResponse{}, err
	}

	defer resp.Body.Close()

	globalEmotesResponse := GlobalEmotesResponse{}
	err = json.NewDecoder(resp.Body).Decode(&globalEmotesResponse)
	if err != nil {
		slog.Error("Error decoding global emotes response", "error", err)
		return GlobalEmotesResponse{}, err
	}

	slog.Info("Got global emotes", "body", globalEmotesResponse)
	return globalEmotesResponse, nil
}

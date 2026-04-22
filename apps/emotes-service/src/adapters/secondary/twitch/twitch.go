package twitch

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"emotes-service/src/environment"

	"golang.org/x/oauth2/clientcredentials"
)

func GetAccessToken() (string, error) {
	oauth2Config := &clientcredentials.Config{
		ClientID:     environment.GetOrFatal("TWITCH_CLIENT_ID"),
		ClientSecret: environment.GetOrFatal("TWITCH_CLIENT_SECRET"),
		TokenURL:     environment.GetOrFatal("TWITCH_OAUTH_ENDPOINT"),
	}

	token, err := oauth2Config.Token(context.Background())
	if err != nil {
		slog.Error("Error getting access token", "error", err)
		return "", err
	}

	return token.AccessToken, nil
}

func GetGlobalEmotes(accessToken string) (any, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/chat/emotes/global", nil)
	if err != nil {
		slog.Error("Error creating global emotes request", "error", err)
		return nil, err
	}

	req.Header.Set("Client-ID", environment.GetOrFatal("TWITCH_CLIENT_ID"))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Error getting global emotes", "error", err)
		return nil, err
	}

	defer resp.Body.Close()

	globalEmotesResponse := map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&globalEmotesResponse)
	if err != nil {
		slog.Error("Error decoding global emotes response", "error", err)
		return nil, err
	}

	slog.Info("Got global emotes", "body", globalEmotesResponse)
	return globalEmotesResponse, nil
}

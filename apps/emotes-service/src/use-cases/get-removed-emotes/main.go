package main

import (
	"context"
	"emotes-service/src/adapters/secondary/projections_store"
	"encoding/json"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

type emoteImagesDTO struct {
	URL1X string `json:"url_1x"`
	URL2X string `json:"url_2x"`
	URL4X string `json:"url_4x"`
}

type emoteDTO struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Format    []string       `json:"format"`
	Scale     []string       `json:"scale"`
	ThemeMode []string       `json:"theme_mode"`
	Images    emoteImagesDTO `json:"images"`
}

type removedEmoteDTO struct {
	Emote     emoteDTO `json:"emote"`
	RemovedAt string   `json:"removedAt"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	channelId := request.PathParameters["channelId"]
	if channelId == "" {
		return jsonResponse(ctx, 400, map[string]string{"message": "channelId is required"})
	}

	items, err := projections_store.QueryRemovedEmotes(ctx, channelId)
	if err != nil {
		slog.ErrorContext(ctx, "QueryRemovedEmotes failed", "channelId", channelId, "error", err)
		return jsonResponse(ctx, 500, map[string]string{"message": "internal error"})
	}

	removedEmotes := make([]removedEmoteDTO, 0, len(items))
	for _, item := range items {
		removedEmotes = append(removedEmotes, removedEmoteDTO{
			Emote: emoteDTO{
				ID:        item.Emote.ID,
				Name:      item.Emote.Name,
				Format:    item.Emote.Format,
				Scale:     item.Emote.Scale,
				ThemeMode: item.Emote.ThemeMode,
				Images: emoteImagesDTO{
					URL1X: item.Emote.Images.URL1X,
					URL2X: item.Emote.Images.URL2X,
					URL4X: item.Emote.Images.URL4X,
				},
			},
			RemovedAt: *item.RemovedAt,
		})
	}

	return jsonResponse(ctx, 200, removedEmotes)
}

func jsonResponse(ctx context.Context, status int, body any) (events.APIGatewayProxyResponse, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		slog.ErrorContext(ctx, "marshall error", "error", err)
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: `{"message":"internal error"}`}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(payload),
	}, nil
}

func main() {
	slog.SetDefault(lambdacontext.NewLogger())
	lambda.Start(handler)
}

package main

import (
	"context"
	"emotes-service/src/adapters/primary/apigw"
	"emotes-service/src/adapters/secondary/projections_store"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

type activeEmoteDTO struct {
	Emote   apigw.Emote `json:"emote"`
	AddedAt string      `json:"addedAt"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	channelId := request.PathParameters["channelId"]
	if channelId == "" {
		return apigw.JSONResponse(ctx, 400, map[string]string{"message": "channelId is required"})
	}

	items, err := projections_store.QueryActiveEmotes(ctx, channelId)
	if err != nil {
		slog.ErrorContext(ctx, "QueryActiveEmotes failed", "channelId", channelId, "error", err)
		return apigw.JSONResponse(ctx, 500, map[string]string{"message": "internal error"})
	}

	activeEmotes := make([]activeEmoteDTO, 0, len(items))
	for _, item := range items {
		if item.Emote == nil {
			continue
		}
		activeEmotes = append(activeEmotes, activeEmoteDTO{
			Emote:   apigw.EmoteFromEvent(item.Emote),
			AddedAt: item.CreatedAt,
		})
	}

	return apigw.JSONResponse(ctx, 200, activeEmotes)
}

func main() {
	slog.SetDefault(lambdacontext.NewLogger())
	lambda.Start(handler)
}

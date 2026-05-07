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

type removedEmoteDTO struct {
	Emote     apigw.Emote `json:"emote"`
	RemovedAt string      `json:"removedAt"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	channelId := request.PathParameters["channelId"]
	if channelId == "" {
		return apigw.JSONResponse(ctx, 400, map[string]string{"message": "channelId is required"})
	}

	items, err := projections_store.QueryRemovedEmotes(ctx, channelId)
	if err != nil {
		slog.ErrorContext(ctx, "QueryRemovedEmotes failed", "channelId", channelId, "error", err)
		return apigw.JSONResponse(ctx, 500, map[string]string{"message": "internal error"})
	}

	removedEmotes := make([]removedEmoteDTO, 0, len(items))
	for _, item := range items {
		removedEmotes = append(removedEmotes, removedEmoteDTO{
			Emote:     apigw.EmoteFromEvent(item.Emote),
			RemovedAt: *item.RemovedAt,
		})
	}

	return apigw.JSONResponse(ctx, 200, removedEmotes)
}

func main() {
	slog.SetDefault(lambdacontext.NewLogger())
	lambda.Start(handler)
}

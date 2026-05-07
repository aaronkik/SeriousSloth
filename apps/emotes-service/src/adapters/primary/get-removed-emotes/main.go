package main

import (
	"context"
	"emotes-service/src/apigw"
	getremovedemotes "emotes-service/src/use-cases/get-removed-emotes"
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

	removedEmotes, err := getremovedemotes.RemovedEmotes(ctx, channelId)

	if err != nil {
		slog.ErrorContext(ctx, "get-removed-emotes use-case failed", "channelId", channelId, "error", err)
		return apigw.JSONResponse(ctx, 500, map[string]string{"message": "internal error"})
	}

	dtos := make([]removedEmoteDTO, 0, len(removedEmotes))
	for _, item := range removedEmotes {
		dtos = append(dtos, removedEmoteDTO{
			Emote:     apigw.EmoteFromEvent(item.Emote),
			RemovedAt: item.RemovedAt,
		})
	}
	return apigw.JSONResponse(ctx, 200, dtos)
}

func main() {
	slog.SetDefault(lambdacontext.NewLogger())
	lambda.Start(handler)
}

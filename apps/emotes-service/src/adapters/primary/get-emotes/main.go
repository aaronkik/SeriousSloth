package main

import (
	"context"
	"emotes-service/src/apigw"
	getemotes "emotes-service/src/use-cases/get-emotes"
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

	activeEmotes, err := getemotes.ActiveEmotes(ctx, channelId)

	if err != nil {
		slog.ErrorContext(ctx, "get-emotes use-case failed", "channelId", channelId, "error", err)
		return apigw.JSONResponse(ctx, 500, map[string]string{"message": "internal error"})
	}

	dtos := make([]activeEmoteDTO, 0, len(activeEmotes))
	for _, item := range activeEmotes {
		dtos = append(dtos, activeEmoteDTO{
			Emote:   apigw.EmoteFromEvent(item.Emote),
			AddedAt: item.AddedAt,
		})
	}
	return apigw.JSONResponse(ctx, 200, dtos)
}

func main() {
	slog.SetDefault(lambdacontext.NewLogger())
	lambda.Start(handler)
}

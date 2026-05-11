package main

import (
	"context"
	"emotes-service/src/adapters/secondary/event_store"
	"emotes-service/src/apigw"
	"emotes-service/src/environment"
	getremovedemotes "emotes-service/src/use-cases/get-removed-emotes"
	"log/slog"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/newrelic/go-agent/v3/integrations/nrlambda"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type removedEmoteDTO struct {
	Emote     apigw.Emote `json:"emote"`
	RemovedAt string      `json:"removedAt"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	channelId := request.PathParameters["channelId"]

	if strings.EqualFold(channelId, event_store.GlobalEmotesAggregateId) {
		channelId = event_store.GlobalEmotesAggregateId
	}

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
	logger := lambdacontext.NewLogger().With(
		slog.Group("tags",
			"project", environment.GetOrFatal("PROJECT"),
			"stack", environment.GetOrFatal("STACK"),
		),
	)
	slog.SetDefault(logger)
	app, err := newrelic.NewApplication(nrlambda.ConfigOption())
	if nil != err {
		slog.Error("error creating app (invalid config)", "error", err)
	}

	nrlambda.Start(handler, app)
}

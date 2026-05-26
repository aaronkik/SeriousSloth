package main

import (
	"context"
	"emotes-service/src/adapters/secondary/event_store"
	"emotes-service/src/apigw"
	"emotes-service/src/environment"
	getemotes "emotes-service/src/use-cases/get-emotes"
	"errors"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/newrelic/go-agent/v3/integrations/nrlambda"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type activeEmoteDTO struct {
	Emote   apigw.Emote `json:"emote"`
	AddedAt string      `json:"addedAt"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	apigw.AnnotateRequest(ctx, request)

	channelId := request.PathParameters["channelId"]

	aggregateId, err := event_store.AggregateIdFromChannelId(channelId)
	if err != nil {
		if errors.Is(err, event_store.ErrInvalidChannelId) {
			return apigw.JSONResponse(ctx, 400, map[string]string{"message": "channelId must be 'global' or a numeric twitch id"})
		}
		slog.ErrorContext(ctx, "aggregate id translation failed", "channelId", channelId, "error", err)
		return apigw.JSONResponse(ctx, 500, map[string]string{"message": "internal error"})
	}

	activeEmotes, err := getemotes.ActiveEmotes(ctx, aggregateId)

	if err != nil {
		slog.ErrorContext(ctx, "get-emotes use-case failed", "channelId", channelId, "aggregateId", aggregateId, "error", err)
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

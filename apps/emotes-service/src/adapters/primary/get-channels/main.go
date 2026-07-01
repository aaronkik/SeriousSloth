package main

import (
	"context"
	"emotes-service/src/apigw"
	"emotes-service/src/environment"
	getchannels "emotes-service/src/use-cases/get-channels"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/newrelic/go-agent/v3/integrations/nrlambda"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	apigw.AnnotateRequest(ctx, request)

	channels, err := getchannels.ListChannels(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "get-channels use-case failed", "error", err)
		return apigw.JSONResponse(ctx, 500, map[string]string{"message": "internal error"})
	}

	dtos := make([]apigw.Channel, 0, len(channels))
	for _, channel := range channels {
		dtos = append(dtos, apigw.Channel{
			Id:          channel.Id,
			TwitchId:    channel.TwitchId,
			DisplayName: channel.DisplayName,
			ImageUrl:    channel.ImageUrl,
			AddedAt:     channel.AddedAt,
			UpdatedAt:   channel.UpdatedAt,
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
	app, err := newrelic.NewApplication(newrelic.ConfigFromEnvironment())
	if nil != err {
		slog.Error("error creating app (invalid config)", "error", err)
	}

	nrlambda.Start(handler, app)
}

package main

import (
	"context"
	"emotes-service/src/apigw"
	"emotes-service/src/environment"
	addchannel "emotes-service/src/use-cases/add-channel"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/newrelic/go-agent/v3/integrations/nrlambda"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type addChannelRequest struct {
	TwitchId    string `json:"twitchId"`
	DisplayName string `json:"displayName"`
	ImageUrl    string `json:"imageUrl"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	apigw.AnnotateRequest(ctx, request)
	txn := newrelic.FromContext(ctx)

	var body addChannelRequest
	if err := json.Unmarshal([]byte(request.Body), &body); err != nil {
		slog.InfoContext(ctx, "add-channel invalid json", "error", err)
		return apigw.JSONResponse(ctx, 400, map[string]string{"message": "invalid request body"})
	}

	if txn != nil && body.TwitchId != "" {
		txn.AddAttribute("channel.twitchId", body.TwitchId)
	}

	channel, err := addchannel.AddChannel(ctx, addchannel.Input{
		TwitchId:    body.TwitchId,
		DisplayName: body.DisplayName,
		ImageUrl:    body.ImageUrl,
	})
	if err != nil {
		if txn != nil {
			txn.NoticeError(err)
		}
		switch {
		case errors.Is(err, addchannel.ErrInvalidInput):
			return apigw.JSONResponse(ctx, 400, map[string]string{"message": "twitchId, displayName and imageUrl are required"})
		case errors.Is(err, addchannel.ErrAlreadyExists):
			return apigw.JSONResponse(ctx, 409, map[string]string{"message": "channel already exists"})
		default:
			slog.ErrorContext(ctx, "add-channel use-case failed", "error", err)
			return apigw.JSONResponse(ctx, 500, map[string]string{"message": "internal error"})
		}
	}

	return apigw.JSONResponse(ctx, 201, apigw.Channel{
		Id:          channel.Id,
		TwitchId:    channel.TwitchId,
		DisplayName: channel.DisplayName,
		ImageUrl:    channel.ImageUrl,
		AddedAt:     channel.AddedAt,
		UpdatedAt:   channel.UpdatedAt,
	})
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

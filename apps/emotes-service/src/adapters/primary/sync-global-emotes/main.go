package main

import (
	"context"
	"emotes-service/src/environment"
	syncglobalemotes "emotes-service/src/use-cases/sync-global-emotes"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/newrelic/go-agent/v3/integrations/nrlambda"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func handler(ctx context.Context, event events.CloudWatchEvent) error {
	slog.InfoContext(ctx, "canonical log", "event", event)
	return syncglobalemotes.Execute(ctx)
}

func main() {
	logger := lambdacontext.NewLogger().With(
		slog.Group("tags",
			"project", environment.GetOrFatal("PROJECT"),
			"stack", environment.GetOrFatal("STACK"),
		),
	)
	slog.SetDefault(logger)

	app, err := newrelic.NewApplication(
		nrlambda.ConfigOption(),
		newrelic.ConfigFromEnvironment(),
	)
	if err != nil {
		slog.Error("error creating app (invalid config)", "err", err)
	}

	nrlambda.Start(handler, app)
}

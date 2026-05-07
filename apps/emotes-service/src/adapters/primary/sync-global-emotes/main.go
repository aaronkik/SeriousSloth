package main

import (
	"context"
	syncglobalemotes "emotes-service/src/use-cases/sync-global-emotes"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

func handler(ctx context.Context, event events.CloudWatchEvent) error {
	slog.InfoContext(ctx, "canonical log", "event", event)
	return syncglobalemotes.Execute(ctx)
}

func main() {
	slog.SetDefault(lambdacontext.NewLogger())
	lambda.Start(handler)
}

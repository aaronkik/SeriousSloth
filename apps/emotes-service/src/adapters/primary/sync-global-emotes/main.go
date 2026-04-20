package main

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

func handler(ctx context.Context, event json.RawMessage) error {
	slog.InfoContext(ctx, "canonical log", "event", event)
	slog.InfoContext(ctx, "hello")
	return nil
}

func main() {
	slog.SetDefault(lambdacontext.NewLogger())
	lambda.Start(handler)
}

package main

import (
	"context"
	"emotes-service/src/environment"
	mocktwitchapi "emotes-service/src/use-cases/mock-twitch-api"
	"log/slog"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/newrelic/go-agent/v3/integrations/nrlambda"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func handler(ctx context.Context, event events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	path := event.RequestContext.HTTP.Path
	key := lookupKey(path, event.RawQueryString)
	slog.InfoContext(ctx, "mock-twitch-api request", "path", path, "rawQuery", event.RawQueryString, "lookupKey", key, "method", event.RequestContext.HTTP.Method)

	body, err := mocktwitchapi.Execute(ctx, key)
	if err != nil {
		return events.LambdaFunctionURLResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       body,
	}, nil
}

// lookupKey builds a deterministic DDB PK from the request path and (optional)
// query string. Backwards compatible: when the request has no query string the
// key is just the path, matching existing seeded fixtures.
func lookupKey(path, rawQuery string) string {
	if rawQuery == "" {
		return path
	}
	values, err := url.ParseQuery(rawQuery)
	if err != nil {
		// Fall back to raw — caller can still seed against the verbatim string.
		return path + "?" + rawQuery
	}
	// url.Values.Encode sorts keys and per-key values, giving a canonical form.
	return path + "?" + values.Encode()
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

package mocktwitchapi

import (
	"context"
	"emotes-service/src/adapters/secondary/mock_store"

	"github.com/newrelic/go-agent/v3/newrelic"
)

func Execute(ctx context.Context, path string) (string, error) {
	txn := newrelic.FromContext(ctx)

	seg := txn.StartSegment("mock_store.Lookup")
	defer seg.End()
	return mock_store.Lookup(ctx, path)
}

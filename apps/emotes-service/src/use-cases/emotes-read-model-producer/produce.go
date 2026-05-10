package readmodelproducer

import (
	"context"
	"emotes-service/src/adapters/secondary/event_store"
	"emotes-service/src/adapters/secondary/projections_store"

	"github.com/newrelic/go-agent/v3/newrelic"
)

func Execute(ctx context.Context, events []event_store.EmoteServiceEvent) error {
	txn := newrelic.FromContext(ctx)

	seg := txn.StartSegment("projections_store.Persist")
	defer seg.End()
	for _, event := range events {
		if err := projections_store.Persist(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

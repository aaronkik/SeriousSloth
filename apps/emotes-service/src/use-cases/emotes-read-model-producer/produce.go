package readmodelproducer

import (
	"context"
	"emotes-service/src/adapters/secondary/event_store"
	"emotes-service/src/adapters/secondary/projections_store"
)

func Execute(ctx context.Context, events []event_store.EmoteServiceEvent) error {
	for _, event := range events {
		if err := projections_store.Persist(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

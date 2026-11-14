package getaddedemotes

import (
	"context"
	"emotes-service/src/adapters/secondary/event_store"
	"log/slog"

	"github.com/newrelic/go-agent/v3/newrelic"
)

type AddedEmote struct {
	Emote   event_store.EmoteServiceEventEmote
	AddedAt string
}

func AddedEmotes(ctx context.Context, aggregateId string) ([]AddedEmote, error) {
	txn := newrelic.FromContext(ctx)

	seg := txn.StartSegment("event_store.LoadEmoteAddedEvents")
	events, err := event_store.LoadEmoteAddedEvents(ctx, aggregateId)
	seg.End()
	if err != nil {
		slog.ErrorContext(ctx, "LoadEmoteAddedEvents failed", "aggregateId", aggregateId, "error", err)
		return nil, err
	}

	results := make([]AddedEmote, 0, len(events))
	for _, event := range events {
		results = append(results, AddedEmote{
			Emote:   *event.Emote,
			AddedAt: event.CreatedAt,
		})
	}
	return results, nil
}

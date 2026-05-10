package getemotes

import (
	"context"
	"emotes-service/src/adapters/secondary/event_store"
	"emotes-service/src/adapters/secondary/projections_store"
	"log/slog"

	"github.com/newrelic/go-agent/v3/newrelic"
)

type ActiveEmote struct {
	Emote   event_store.EmoteServiceEventEmote
	AddedAt string
}

func ActiveEmotes(ctx context.Context, channelId string) ([]ActiveEmote, error) {
	txn := newrelic.FromContext(ctx)

	seg := txn.StartSegment("projections_store.QueryActiveEmotes")
	items, err := projections_store.QueryActiveEmotes(ctx, channelId)
	seg.End()
	if err != nil {
		slog.ErrorContext(ctx, "QueryActiveEmotes failed", "channelId", channelId, "error", err)
		return nil, err
	}

	results := make([]ActiveEmote, 0, len(items))
	for _, item := range items {
		results = append(results, ActiveEmote{
			Emote:   item.Emote,
			AddedAt: item.CreatedAt,
		})
	}
	return results, nil
}

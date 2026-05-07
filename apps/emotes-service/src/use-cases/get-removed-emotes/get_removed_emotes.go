package getremovedemotes

import (
	"context"
	"emotes-service/src/adapters/secondary/event_store"
	"emotes-service/src/adapters/secondary/projections_store"
	"log/slog"
)

type RemovedEmote struct {
	Emote     event_store.EmoteServiceEventEmote
	RemovedAt string
}

func RemovedEmotes(ctx context.Context, channelId string) ([]RemovedEmote, error) {
	items, err := projections_store.QueryRemovedEmotes(ctx, channelId)
	if err != nil {
		slog.ErrorContext(ctx, "QueryRemovedEmotes failed", "channelId", channelId, "error", err)
		return nil, err
	}

	results := make([]RemovedEmote, 0, len(items))
	for _, item := range items {
		results = append(results, RemovedEmote{
			Emote:     item.Emote,
			RemovedAt: *item.RemovedAt,
		})
	}
	return results, nil
}

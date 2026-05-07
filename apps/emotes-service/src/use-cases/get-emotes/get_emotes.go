package getemotes

import (
	"context"
	"emotes-service/src/adapters/secondary/event_store"
	"emotes-service/src/adapters/secondary/projections_store"
	"log/slog"
)

type ActiveEmote struct {
	Emote   event_store.EmoteServiceEventEmote
	AddedAt string
}

func ActiveEmotes(ctx context.Context, channelId string) ([]ActiveEmote, error) {
	items, err := projections_store.QueryActiveEmotes(ctx, channelId)
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

package getchannels

import (
	"context"
	"emotes-service/src/adapters/secondary/channels_store"
	"log/slog"
)

type Channel struct {
	Id          string
	TwitchId    string
	DisplayName string
	ImageUrl    string
	AddedAt     string
	UpdatedAt   string
}

func ListChannels(ctx context.Context) ([]Channel, error) {
	items, err := channels_store.QueryAll(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "channels_store.QueryAll failed", "error", err)
		return nil, err
	}

	results := make([]Channel, 0, len(items))
	for _, item := range items {
		results = append(results, Channel{
			Id:          item.Id,
			TwitchId:    item.TwitchId,
			DisplayName: item.DisplayName,
			ImageUrl:    item.ImageUrl,
			AddedAt:     item.AddedAt,
			UpdatedAt:   item.UpdatedAt,
		})
	}
	return results, nil
}

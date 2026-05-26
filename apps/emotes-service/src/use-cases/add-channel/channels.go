package addchannel

import (
	"context"
	"emotes-service/src/adapters/secondary/channels_store"
	"emotes-service/src/adapters/secondary/twitch"
	"emotes-service/src/ids"
	"errors"
	"log/slog"
	"strings"
	"time"
)

var (
	ErrInvalidInput    = errors.New("invalid channel input")
	ErrChannelNotFound = errors.New("twitch channel not found")
	ErrAlreadyExists   = channels_store.ErrAlreadyExists
)

type Input struct {
	TwitchId string
}

type Channel struct {
	Id          string
	TwitchId    string
	DisplayName string
	ImageUrl    string
	AddedAt     string
	UpdatedAt   string
}

func AddChannel(ctx context.Context, input Input) (Channel, error) {
	twitchId := strings.TrimSpace(input.TwitchId)

	if twitchId == "" {
		return Channel{}, ErrInvalidInput
	}

	accessToken, err := twitch.GetAccessToken(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "twitch.GetAccessToken failed", "twitchId", twitchId, "error", err)
		return Channel{}, err
	}

	user, err := twitch.GetUserById(ctx, accessToken, twitchId)
	if err != nil {
		slog.ErrorContext(ctx, "twitch.GetUserById failed", "twitchId", twitchId, "error", err)
		return Channel{}, err
	}
	if user == nil {
		return Channel{}, ErrChannelNotFound
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)
	item := channels_store.ChannelItem{
		PK:          channels_store.ChannelsPartitionKey,
		SK:          twitchId,
		Id:          ids.New("chnl_"),
		TwitchId:    twitchId,
		DisplayName: user.DisplayName,
		ImageUrl:    user.ProfileImageURL,
		AddedAt:     now,
		UpdatedAt:   now,
	}

	if err := channels_store.Put(ctx, item); err != nil {
		if errors.Is(err, channels_store.ErrAlreadyExists) {
			return Channel{}, ErrAlreadyExists
		}
		slog.ErrorContext(ctx, "channels_store.Put failed", "twitchId", twitchId, "error", err)
		return Channel{}, err
	}

	return Channel{
		Id:          item.Id,
		TwitchId:    item.TwitchId,
		DisplayName: item.DisplayName,
		ImageUrl:    item.ImageUrl,
		AddedAt:     item.AddedAt,
		UpdatedAt:   item.UpdatedAt,
	}, nil
}

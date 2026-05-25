package addchannel

import (
	"context"
	"crypto/rand"
	"emotes-service/src/adapters/secondary/channels_store"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"time"
)

var (
	ErrInvalidInput  = errors.New("invalid channel input")
	ErrAlreadyExists = channels_store.ErrAlreadyExists
)

type Input struct {
	TwitchId    string
	DisplayName string
	ImageUrl    string
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
	displayName := strings.TrimSpace(input.DisplayName)
	imageUrl := strings.TrimSpace(input.ImageUrl)

	if twitchId == "" || displayName == "" || imageUrl == "" {
		return Channel{}, ErrInvalidInput
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)
	item := channels_store.ChannelItem{
		PK:          channels_store.ChannelsPartitionKey,
		SK:          twitchId,
		Id:          generateId(),
		TwitchId:    twitchId,
		DisplayName: displayName,
		ImageUrl:    imageUrl,
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

func generateId() string {
	bytes := make([]byte, 12)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("chnl_%s", hex.EncodeToString(bytes))
}

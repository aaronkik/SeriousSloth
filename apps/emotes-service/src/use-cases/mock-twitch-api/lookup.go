package mocktwitchapi

import (
	"context"
	"emotes-service/src/adapters/secondary/mock_store"
)

func Execute(ctx context.Context, path string) (string, error) {
	return mock_store.Lookup(ctx, path)
}

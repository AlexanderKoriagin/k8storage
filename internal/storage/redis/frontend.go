package redis

import (
	"context"
	"fmt"

	"github.com/akrillis/k8storage/internal/entities"
	"github.com/akrillis/k8storage/internal/entities/config"
	"github.com/go-redis/redis/v8"
)

// Frontend - redis connection for frontend server
type Frontend struct {
	client *redis.Client
}

func NewFrontend(ctx context.Context, config *config.Redis) (*Frontend, error) {
	client, err := connect(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("couldn't create redis client: %w", err)
	}

	return &Frontend{client: client}, nil
}

// Put - put file to short term storage.
func (f *Frontend) Put(ctx context.Context, req *entities.PutFileRequest) error {
	key := fmt.Sprintf(keyClientData, req.ClientID, req.Name)
	if err := f.client.Set(ctx, key, req.Data, 0).Err(); err != nil && err != Nil {
		return fmt.Errorf("couldn't set key %s: %w", key, err)
	}

	if err := f.client.LPush(ctx, requestList, key).Err(); err != nil && err != Nil {
		_ = f.client.Del(ctx, key).Err()
		return fmt.Errorf("couldn't push key %s to list %s: %w", key, requestList, err)
	}

	return nil
}

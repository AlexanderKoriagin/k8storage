package redis

import (
	"context"
	"fmt"

	"github.com/akrillis/k8storage/internal/entities/config"
	"github.com/go-redis/redis/v8"
)

const (
	Nil = redis.Nil

	keyClientData = "%s,%s"
	requestList   = "clientsRequests"
)

func connect(ctx context.Context, config *config.Redis) (*redis.Client, error) {
	rc := redis.NewClient(
		&redis.Options{
			Addr:         config.Address,
			Password:     config.Password,
			DB:           config.DB,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
		})

	if rc == nil {
		return nil, fmt.Errorf("could not create redis client")
	}

	if err := rc.Ping(ctx).Err(); err != nil && err != Nil {
		return nil, fmt.Errorf("could not ping redis: %w", err)
	}

	return rc, nil
}

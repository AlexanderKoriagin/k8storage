package storage

import (
	"context"

	"github.com/akrillis/k8storage/internal/entities"
)

//go:generate mockgen -source=storage.go -package=storage -destination=storage_mock.go

// Frontender is an interface for short term storage frontend.
type Frontender interface {
	Put(ctx context.Context, req *entities.PutFileRequest) error
}

// Backender is an interface for long term storage backend.
type Backender interface {
	Get(ctx context.Context, req *entities.GetFileRequest) ([]byte, error)
}

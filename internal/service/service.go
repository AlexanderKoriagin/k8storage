package service

import (
	"context"

	"github.com/akrillis/k8storage/internal/entities"
)

//go:generate mockgen -source=service.go -package=service -destination=service_mock.go

// Receiver is an interface for putting data to short term files.
type Receiver interface {
	Put(ctx context.Context, req *entities.PutFileRequest) error
}

// Restorer is an interface for getting data from long term files.
type Restorer interface {
	Get(ctx context.Context, req *entities.GetFileRequest) ([]byte, error)
}

package backend

import (
	"context"

	"github.com/akrillis/k8storage/internal/entities"
	"github.com/akrillis/k8storage/internal/storage"
)

type Restorer struct {
	sr storage.Backender
}

func New(sr storage.Backender) *Restorer {
	return &Restorer{
		sr: sr,
	}
}

// Get returns file data by GetFileRequest parameters.
func (r *Restorer) Get(ctx context.Context, req *entities.GetFileRequest) ([]byte, error) {
	return r.sr.Get(ctx, req)
}

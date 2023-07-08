package receiver

import (
	"context"

	"github.com/akrillis/k8storage/internal/entities"
	"github.com/akrillis/k8storage/internal/storage"
)

type Receiver struct {
	sf storage.Frontender
}

func New(sf storage.Frontender) *Receiver {
	return &Receiver{
		sf: sf,
	}
}

// Put stores the file in the short term storage.
func (r *Receiver) Put(ctx context.Context, req *entities.PutFileRequest) error {
	return r.sf.Put(ctx, req)
}

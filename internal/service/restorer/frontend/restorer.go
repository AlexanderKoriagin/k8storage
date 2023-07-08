package frontend

import (
	"context"
	"fmt"

	"github.com/akrillis/k8storage/internal/entities"
	grpcBackend "github.com/akrillis/k8storage/pkg/grpc/backend"
)

type Restorer struct {
	bc grpcBackend.BackendClient
}

func New(bc grpcBackend.BackendClient) *Restorer {
	return &Restorer{
		bc: bc,
	}
}

// Get returns file data by GetFileRequest parameters.
func (r *Restorer) Get(ctx context.Context, req *entities.GetFileRequest) ([]byte, error) {
	data, err := r.bc.GetFile(
		ctx,
		&grpcBackend.GetRequest{
			ClientId: req.ClientID,
			Name:     req.Name,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	return data.Data, nil
}

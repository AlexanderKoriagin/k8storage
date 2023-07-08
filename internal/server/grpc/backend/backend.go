package backend

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/akrillis/k8storage/internal/entities"
	"github.com/akrillis/k8storage/internal/entities/config"
	"github.com/akrillis/k8storage/internal/service"
	grpcBackend "github.com/akrillis/k8storage/pkg/grpc/backend"
	"google.golang.org/grpc"
)

type server struct {
	service service.Restorer
	*grpcBackend.UnimplementedBackendServer
}

// ServeGRPC creates a new configured grpcBackend.BackendServer.
func ServeGRPC(cfg *config.Grpc) error {
	lis, err := net.Listen("tcp", cfg.Listen)
	if err != nil {
		return fmt.Errorf("failed to listen on %q: %w", cfg.Listen, err)
	}

	srv := grpc.NewServer()
	s := &server{
		service:                    cfg.Restorer,
		UnimplementedBackendServer: &grpcBackend.UnimplementedBackendServer{},
	}

	grpcBackend.RegisterBackendServer(srv, s)

	// Start serving requests.
	e := make(chan error, 1)
	go func() {
		e <- srv.Serve(lis)
	}()

	// Wait for shutdown signal for graceful shutdown.
	select {
	case <-cfg.Ctx.Done():
		t := time.AfterFunc(cfg.GracePeriod, srv.Stop)
		srv.GracefulStop()
		t.Stop()
		cfg.Wg.Done()
		return nil
	case err := <-e:
		return err
	}
}

func (s *server) GetFile(ctx context.Context, r *grpcBackend.GetRequest) (*grpcBackend.GetResponse, error) {
	data, err := s.service.Get(
		ctx,
		&entities.GetFileRequest{
			ClientID: r.GetClientId(),
			Name:     r.GetName(),
		},
	)
	if err != nil {
		return nil, err
	}

	return &grpcBackend.GetResponse{
		Data: data,
	}, nil
}

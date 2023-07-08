package config

import (
	"context"
	"sync"
	"time"

	"github.com/akrillis/k8storage/internal/service"
)

// Grpc is a configuration for gRPC server.
type Grpc struct {
	Ctx         context.Context  // context for graceful shutdown
	Listen      string           // address to listen
	Restorer    service.Restorer // service to handle restore requests
	GracePeriod time.Duration    // grace period for graceful shutdown
	Wg          *sync.WaitGroup
}

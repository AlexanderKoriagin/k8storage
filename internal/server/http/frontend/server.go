package frontend

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/akrillis/k8storage/internal/entities/config"
)

// StartServer starts http server and watch for context cancellation for graceful shutdown.
func StartServer(cfg *config.Http) error {
	server := &http.Server{
		Addr:              fmt.Sprintf("0.0.0.0:%d", cfg.Listen),
		Handler:           cfg.Router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
	}

	go func(ctx context.Context, gracePeriod time.Duration, wg *sync.WaitGroup) {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), gracePeriod)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("failed to graceful shutdown server: %v", err)
		}
		wg.Done()
	}(cfg.Ctx, cfg.GracePeriod, cfg.Wg)

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to listen and serve HTTP requests: %w", err)
	}

	return nil
}

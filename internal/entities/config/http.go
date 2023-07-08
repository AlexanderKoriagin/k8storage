package config

import (
	"context"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// Http is a configuration for http server.
type Http struct {
	Ctx         context.Context
	Listen      uint16
	Router      *mux.Router
	GracePeriod time.Duration
	Wg          *sync.WaitGroup
}

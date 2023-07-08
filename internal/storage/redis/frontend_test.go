package redis

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/akrillis/k8storage/internal/entities"
	"github.com/akrillis/k8storage/internal/entities/config"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFrontend_Put(t *testing.T) {
	s, err := miniredis.Run()
	require.NoError(t, err)
	defer s.Close()

	cfg := &config.Redis{
		Address:      s.Addr(),
		Password:     "",
		DB:           0,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	f, err := NewFrontend(context.Background(), cfg)
	require.NoError(t, err)

	req := &entities.PutFileRequest{
		ClientID: "0",
		Name:     "1.json",
		Data:     []byte(`{"test": "test"}`),
	}

	// Context timeout
	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 10*time.Nanosecond)
	defer cancelTimeout()

	assert.Error(t, f.Put(ctxTimeout, req))

	// Success
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	assert.NoError(t, f.Put(ctx, req))

	rc, err := connect(ctx, cfg)
	require.NoError(t, err)
	defer rc.Close()

	err = rc.Get(ctx, fmt.Sprintf(keyClientData, req.ClientID, req.Name)).Err()
	require.NoError(t, err)

	require.Equal(t, int64(1), rc.LLen(ctx, requestList).Val())
}

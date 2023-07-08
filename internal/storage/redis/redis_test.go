package redis

import (
	"context"
	"testing"
	"time"

	"github.com/akrillis/k8storage/internal/entities/config"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/require"
)

func Test_connect(t *testing.T) {
	s, err := miniredis.Run()
	require.NoError(t, err)
	defer s.Close()

	// Wrong address
	cfg := &config.Redis{
		Address:      "0.0.0.0:12345",
		Password:     "",
		DB:           0,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	c, err := connect(context.Background(), cfg)
	require.Error(t, err)
	require.Nil(t, c)

	// Context timeout
	cfg.Address = s.Addr()
	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 10*time.Nanosecond)
	defer cancelTimeout()

	c, err = connect(ctxTimeout, cfg)
	require.Error(t, err)
	require.Nil(t, c)

	// Correct address
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	c, err = connect(ctx, cfg)
	require.NoError(t, err)
	require.NotNil(t, c)
}

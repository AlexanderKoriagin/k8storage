package redis

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/akrillis/k8storage/internal/entities"
	"github.com/akrillis/k8storage/internal/entities/config"
	"github.com/go-redis/redis/v8"
)

const (
	saverTimeout = 100 * time.Millisecond
)

// Backend - redis connection for backend server
type Backend struct {
	clients  clients
	chanStop chan struct{}
	wg       *sync.WaitGroup
	lwg      *sync.WaitGroup
}

type clients struct {
	source  *redis.Client
	storage []*redis.Client
}

func NewBackend(ctx context.Context, cfg *config.StorageBackend) (*Backend, error) {
	var (
		storage = make([]*redis.Client, len(cfg.StorageCfg))
		err     error
	)

	source, err := connect(ctx, cfg.SourceCfg)
	if err != nil {
		return nil, fmt.Errorf("couldn't create redis client for server %s: %w", cfg.SourceCfg.Address, err)
	}

	for i := range cfg.StorageCfg {
		storage[i], err = connect(ctx, cfg.StorageCfg[i])
		if err != nil {
			return nil, fmt.Errorf("couldn't create redis client for server %s: %w", cfg.StorageCfg[i].Address, err)
		}
	}

	b := &Backend{
		clients: clients{
			source:  source,
			storage: storage,
		},
		chanStop: cfg.ChanStop,
		wg:       cfg.Wg,
		lwg:      new(sync.WaitGroup),
	}

	go b.saver()

	return b, nil
}

// Get - get file content from long term storage.
func (b *Backend) Get(ctx context.Context, req *entities.GetFileRequest) ([]byte, error) {
	var (
		key = fmt.Sprintf(keyClientData, req.ClientID, req.Name)
		out []byte
	)

	for i := range b.clients.storage {
		data, err := b.clients.storage[i].Get(ctx, key).Bytes()
		if err != nil {
			if err == Nil {
				break
			}
			return nil, fmt.Errorf("couldn't get data for key %s for storage %d: %w", key, i, err)
		}

		out = append(out, data...)
	}

	return out, nil
}

// saver is a backend worker that saves data from short term storage to long term storage.
//
//	it runs in a separate goroutine
//	it is stopped by closing chanStop channel
//	every saverTimeout it checks tasks pipe for new tasks
//	it gets data from short term storage according task and saves data into long term storage
func (b *Backend) saver() {
	for {
		select {
		case <-b.chanStop:
			b.lwg.Wait() // wait for all saver goroutines to finish

			_ = b.clients.source.Close()
			for i := range b.clients.storage {
				_ = b.clients.storage[i].Close()
			}

			b.wg.Done()
			return
		case <-time.After(saverTimeout):
			b.lwg.Add(1)
			go func() {
				defer b.lwg.Done()

				key, data, err := b.getSourceData()
				if err != nil {
					log.Printf("couldn't get source data: %v", err)
					return
				}

				if len(data) == 0 {
					return
				}

				if err = b.saveSourceData(key, data); err != nil {
					log.Printf("couldn't save source data for key %s: %v", key, err)
				}
			}()
		}
	}
}

// getSourceData gets data from short term storage.
//
//	gets task from tasks pipe
//	gets data from short term storage according task
func (b *Backend) getSourceData() (string, []byte, error) {
	var (
		key  string
		data []byte
		err  error
	)
	key, err = b.clients.source.RPop(context.TODO(), requestList).Result()
	if err != nil && err != Nil {
		return key, data, fmt.Errorf("couldn't get key from list %s: %s", requestList, err)
	}

	data, err = b.clients.source.Get(context.TODO(), key).Bytes()
	if err != nil && err != Nil {
		return key, data, fmt.Errorf("couldn't get data from key %s: %s", key, err)
	}

	return key, data, nil
}

// saveSourceData saves data to long term storage.
//
//	divide data into chunks
//	save each chunk to separate storage (round robin mechanism)
func (b *Backend) saveSourceData(key string, data []byte) error {
	var chunkSize = len(data)/len(b.clients.storage) - 1
	for i := 0; i < len(b.clients.storage); i++ {
		var chunk []byte
		switch i {
		case len(b.clients.storage) - 1:
			chunk = data[i*chunkSize:]
		default:
			chunk = data[i*chunkSize : (i+1)*chunkSize]
		}

		if err := b.clients.storage[i].Set(context.TODO(), key, chunk, 0).Err(); err != nil && err != Nil {
			return fmt.Errorf("couldn't set data to key %s for storage %d: %s", key, i, err)
		}
	}

	return nil
}

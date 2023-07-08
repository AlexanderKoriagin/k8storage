package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/akrillis/k8storage/internal/entities/config"
	"github.com/akrillis/k8storage/internal/server/grpc/backend"
	svcRestorerBackend "github.com/akrillis/k8storage/internal/service/restorer/backend"
	storage "github.com/akrillis/k8storage/internal/storage/redis"
	"github.com/jessevdk/go-flags"
)

type Opts struct {
	SourceRedisAddress               string `long:"sourceredisaddress" description:"source redis node:port address: --redisaddress 1.2.3.4:1234" required:"true"`
	SourceRedisAuth                  string `long:"sourceredisauth" description:"source redis auth" required:"true"`
	SourceRedisDB                    int    `long:"sourceredisdb" default:"0" description:"source redis db number to store data: 0-15"`
	SourceRedisTimeoutReadInSeconds  uint32 `long:"sourceredistimeoutreadinseconds" default:"60" description:"source redis read timeout in seconds"`
	SourceRedisTimeoutWriteInSeconds uint32 `long:"sourceredistimeoutwriteinseconds" default:"60" description:"source redis write timeout in seconds"`

	StorageRedisAddress               []string `long:"storageredisaddress" description:"storage redis node:port address: --redisaddress 1.2.3.4:1234" required:"true"`
	StorageRedisAuth                  []string `long:"storageredisauth" description:"storage redis auth" required:"true"`
	StorageRedisDB                    int      `long:"storageredisdb" default:"0" description:"storage redis db number to store data: 0-15"`
	StorageRedisTimeoutReadInSeconds  uint32   `long:"storageredistimeoutreadinseconds" default:"60" description:"storage redis read timeout in seconds"`
	StorageRedisTimeoutWriteInSeconds uint32   `long:"storageredistimeoutwriteinseconds" default:"60" description:"storage redis write timeout in seconds"`

	GRPCListen string `long:"grpclisten" description:"grpc listen address: -grpclisten :12345" default:":58081"`
}

const (
	shutdownTimeout = 15 * time.Second
)

var (
	opts     Opts
	wg       = new(sync.WaitGroup)
	sigs     = make(chan os.Signal, 1)
	chanStop = make(chan struct{}, 1)
)

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		fmt.Println("Parse flags", err)
		os.Exit(1)
	}

	if err := validateStorageConfig(); err != nil {
		log.Fatalf("Validate storage config: %v", err)
	}

	// prepare configs for source and target storage
	var (
		sourceCfg = &config.Redis{
			Address:      opts.SourceRedisAddress,
			Password:     opts.SourceRedisAuth,
			DB:           opts.SourceRedisDB,
			ReadTimeout:  time.Duration(opts.SourceRedisTimeoutReadInSeconds) * time.Second,
			WriteTimeout: time.Duration(opts.SourceRedisTimeoutWriteInSeconds) * time.Second,
		}
		storageCfg []*config.Redis
	)
	for i := range opts.StorageRedisAddress {
		storageCfg = append(
			storageCfg,
			&config.Redis{
				Address:      opts.StorageRedisAddress[i],
				Password:     opts.StorageRedisAuth[i],
				DB:           opts.StorageRedisDB,
				ReadTimeout:  time.Duration(opts.StorageRedisTimeoutReadInSeconds) * time.Second,
				WriteTimeout: time.Duration(opts.StorageRedisTimeoutWriteInSeconds) * time.Second,
			},
		)
	}

	ctxConnect, cancelConnect := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelConnect()

	// initialize source and target storages
	wg.Add(1)
	sr, err := storage.NewBackend(
		ctxConnect,
		&config.StorageBackend{
			SourceCfg:  sourceCfg,
			StorageCfg: storageCfg,
			ChanStop:   chanStop,
			Wg:         wg,
		},
	)
	if err != nil {
		log.Fatalf("couldn't create storage backend: %v", err)
	}

	// start grpc server for frontend requests
	ctxApp, cancelApp := context.WithCancel(context.Background())
	wg.Add(1)
	go func() {
		if err := backend.ServeGRPC(
			&config.Grpc{
				Ctx:         ctxApp,
				Listen:      opts.GRPCListen,
				Restorer:    svcRestorerBackend.New(sr),
				GracePeriod: shutdownTimeout,
				Wg:          wg,
			},
		); err != nil {
			log.Fatalf("couldn't create grpc server: %v", err)
		}
	}()

	log.Println("backend app started")
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigs
	cancelApp()
	chanStop <- struct{}{}
	wg.Wait()
	log.Println("frontend app stopped")
}

// validateStorageConfig validates long term storage config
//
//	length of storageredisaddress and storageredisauth must be equal
func validateStorageConfig() error {
	if len(opts.StorageRedisAddress) != len(opts.StorageRedisAuth) {
		return fmt.Errorf("storage redis config: length of storageredisaddress and storageredisauth must be equal")
	}

	return nil
}

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
	"github.com/akrillis/k8storage/internal/server/http/frontend"
	httpFrontend "github.com/akrillis/k8storage/internal/server/http/frontend"
	"github.com/akrillis/k8storage/internal/service/receiver"
	restorerFrontend "github.com/akrillis/k8storage/internal/service/restorer/frontend"
	storage "github.com/akrillis/k8storage/internal/storage/redis"
	grpcBackend "github.com/akrillis/k8storage/pkg/grpc/backend"
	"github.com/jessevdk/go-flags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Opts struct {
	RedisAddress               string `long:"redisaddress" description:"redis node:port address: --redisaddress 1.2.3.4:1234" required:"true"`
	RedisAuth                  string `long:"redisauth" description:"redis auth" required:"true"`
	RedisDB                    int    `long:"redisdb" default:"0" description:"redis db number to store data: 0-15"`
	RedisTimeoutReadInSeconds  uint32 `long:"redistimeoutreadinseconds" default:"60" description:"redis read timeout in seconds"`
	RedisTimeoutWriteInSeconds uint32 `long:"redistimeoutwriteinseconds" default:"60" description:"redis write timeout in seconds"`

	GrpcBackendAddress string `long:"grpcbackendaddress" description:"grpc backend address: --grpcbackendaddress 1.2.3.4:1234" required:"true"`

	HTTPListen uint16 `long:"httplisten" description:"http listen port" default:"58080"`
}

const (
	shutdownTimeout = 15 * time.Second
)

var (
	opts Opts
	wg   = new(sync.WaitGroup)
	sigs = make(chan os.Signal, 1)
)

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		fmt.Println("Parse flags", err)
		os.Exit(1)
	}

	ctxConnect, cancelConnect := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelConnect()

	// initialize connect to short time storage
	fr, err := storage.NewFrontend(
		ctxConnect,
		&config.Redis{
			Address:      opts.RedisAddress,
			Password:     opts.RedisAuth,
			DB:           opts.RedisDB,
			ReadTimeout:  time.Duration(opts.RedisTimeoutReadInSeconds) * time.Second,
			WriteTimeout: time.Duration(opts.RedisTimeoutWriteInSeconds) * time.Second,
		})
	if err != nil {
		log.Fatalf("couldn't open redis connection: %v", err)
	}
	rc := receiver.New(fr)

	// initialize connect to backend grpc server to get file content
	conn, err := grpc.Dial(opts.GrpcBackendAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("couldn't open grpc connection: %v", err)
	}
	rs := restorerFrontend.New(grpcBackend.NewBackendClient(conn))

	// start http server to serve user requests
	ctxApp, cancelApp := context.WithCancel(context.Background())
	wg.Add(1)
	go func() {
		if err := httpFrontend.StartServer(
			&config.Http{
				Ctx:         ctxApp,
				Listen:      opts.HTTPListen,
				Router:      frontend.NewRoutes(rc, rs),
				GracePeriod: shutdownTimeout,
				Wg:          wg,
			}); err != nil {
			fmt.Println("couldn't start http server", err)
			os.Exit(1)
		}
	}()

	log.Println("frontend app started")
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigs
	cancelApp()
	wg.Wait()

	if err = conn.Close(); err != nil {
		log.Println("couldn't close grpc connection", err)
	}
	log.Println("frontend app stopped")
}

package main

import (
	"context"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"grumpy-console-companion/sotle-go/config"
	"grumpy-console-companion/sotle-go/internal/server"
	"grumpy-console-companion/sotle-go/pkg/logging"
	"grumpy-console-companion/sotle-go/pkg/storage/mongodb"
)

var (
	env  string
	addr = flag.String("addr", "127.0.0.1:80", "The address to run on.")
)

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf, err := config.LoadConf(env)
	if err != nil {
		panic(fmt.Sprintf("error load config: %s", err.Error()))
	}

	logger, err := logging.NewLogger(conf.Logging)
	if err != nil {
		panic(fmt.Sprintf("Can't initialize logger: %s", err.Error()))
	}

	storage, err := mongodb.New(ctx, conf, logger)
	if err != nil {
		logger.Fatal("error", zap.Error(err))
	}

	s, err := server.New(*addr, storage)
	if err != nil {
		panic(err)
	}
	done := make(chan error, 1)

	logger.Info(fmt.Sprintf("Starting server at: %s", *addr))
	go func() {
		defer close(done)
		done <- s.Start()
	}()

	err = <-done
	logger.Info(fmt.Sprintf("Server exited with error: %s", err))
}

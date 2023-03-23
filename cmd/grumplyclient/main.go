package main

import (
	"context"
	"flag"
	"fmt"
	"grumpy-console-companion/sotle-go/config"
	"grumpy-console-companion/sotle-go/pkg/grcpclient"
	"grumpy-console-companion/sotle-go/pkg/grumpy"
	"grumpy-console-companion/sotle-go/pkg/logging"
)

var (
	env string
)

func init() {
	flag.StringVar(&env, "env", "dev", `Set's run environment. Possible values are "dev" and "prod"`)
	flag.Parse()
}

func main() {
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

	grcpClient, err := grcpclient.New(conf.API.Address)
	if err != nil {
		panic(err)
	}

	brain := grumpy.New(logger, grcpClient)
	brain.Run(ctx)
}

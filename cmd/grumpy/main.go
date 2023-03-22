package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"grumpy-console-companion/sotle-go/config"
	"grumpy-console-companion/sotle-go/pkg/client"
	"grumpy-console-companion/sotle-go/pkg/grumpy/brain/answer"
	"grumpy-console-companion/sotle-go/pkg/logging"
	"grumpy-console-companion/sotle-go/pkg/storage/mongodb"
	"math/rand"
	"time"
)

var (
	env     string
	morning = make(map[int]string)
)

func init() {
	morning[1] = "Yep! Good morning. What for is breakfast!?"
	morning[2] = "Give me the Carrots!"
	morning[3] = "A bit of music?"
	morning[4] = "Let's drink ourself to oblivion!?"

	flag.StringVar(&env, "env", "dev", `Set's run environment. Possible values are "dev" and "prod"`)
	flag.Parse()
}

type ListeningContract interface {
	WaitingForAnswer() string
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

	//storage, err := mongodb.InitDatabase(ctx, conf, logger)
	//if err != nil {
	//	logger.Fatal("error", zap.Error(err))
	//}
	var storage mongodb.DB

	c, err := client.New(conf.API.Address)
	if err != nil {
		panic(err)
	}

	brain := answer.New(&storage)
	run(ctx, brain, logger, c)
}

func run(ctx context.Context, brain ListeningContract, logger *zap.Logger, c *client.Client) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		topic, thought, err := c.QOTD(context.Background(), "greeting")
		if err != nil {
			panic(err)
		}
		err = WriteText(topic)
		if err != nil {
			panic(err)
		}

		err = WriteText(thought)
		if err != nil {
			panic(err)
		}

		aws := brain.WaitingForAnswer()
		err = WriteText(aws)
		if err != nil {
			panic(err)
		}

		time.Sleep(randomSleep())
	}
}

func randomSleep() time.Duration {
	rand.NewSource(time.Now().UnixNano())
	idMin := 10
	idMax := 90

	n := rand.Intn(idMax-idMin) + idMin
	timeToThink := time.Duration(n) * time.Second
	return timeToThink
}

func WriteText(text string) error {
	amount := len(text)

	if text != "" {
		for i := 0; i < amount; i++ {
			fmt.Printf("%v", text[i:i+1])
			time.Sleep(70 * time.Millisecond)
		}
		fmt.Printf("\n")
		return nil
	}

	return errors.New("should not be empty")
}

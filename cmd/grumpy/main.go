package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"grumpy-console-companion/sotle-go/src/config"
	"grumpy-console-companion/sotle-go/src/grumpy/brain/answer"
	"grumpy-console-companion/sotle-go/src/logging"
	"grumpy-console-companion/sotle-go/src/storage/mongodb"
	"math/rand"
	"time"
)

var (
	env     string
	morning = make(map[int]string)
)

func init() {
	morning[1] = "Yep! Good morning. What for breakfast!?"
	morning[2] = "Give me the coffee!"
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

	storage, err := mongodb.InitDatabase(ctx, conf, logger)

	brain := answer.New(storage)
	run(ctx, brain)
}

func run(ctx context.Context, brain ListeningContract) {

	rand.NewSource(time.Now().UnixNano())
	idMin := 1
	idMax := 4

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		goToBrain := rand.Intn(idMax-idMin) + idMin
		thought := morning[goToBrain]
		err := WriteText(thought)
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

package grumpy

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"grumpy-console-companion/sotle-go/pkg/grcpclient"
	"grumpy-console-companion/sotle-go/pkg/grumpy/brain/answer"
	"grumpy-console-companion/sotle-go/pkg/grumpy/dictionary"
	"math/rand"
	"time"
)

type Events int

const (
	WakeUp Events = iota
	Sleep
)

type Brain struct {
	Listening  ListeningContract
	Event      chan Event
	logger     *zap.Logger
	dictionary *dictionary.Dictionary
}

type Event struct {
	Type    Events
	Message string
}

type ListeningContract interface {
	WaitingForAnswer() string
}

func New(logger *zap.Logger, grcp *grcpclient.Client) *Brain {
	return &Brain{
		logger:     logger,
		Listening:  answer.New(),
		Event:      make(chan Event, 1),
		dictionary: dictionary.New(grcp),
	}
}

func (b *Brain) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		//case e := <-b.Event:

		default:
		}

		topic := b.dictionary.Greetings.GetRandom()

		err := WriteText(topic.Phrase)
		if err != nil {
			panic(err)
		}

		aws := b.Listening.WaitingForAnswer()
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

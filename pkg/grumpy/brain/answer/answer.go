package answer

import (
	"bufio"
	"context"
	"grumpy-console-companion/sotle-go/pkg/storage/mongodb"
	"os"
	"time"
)

type respond struct {
	text string
	ok   bool
}

type Listening struct {
	Answer        string
	AnswerChannel chan respond
}

func New(db *mongodb.DB) *Listening {
	return &Listening{
		AnswerChannel: make(chan respond, 1),
	}
}

func (a *Listening) WaitingForAnswer() string {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()

	rsp := a.waiting(ctx)
	if rsp.ok {
		return rsp.text
	}
	return "Oh yeah. Just ignore me! As Usual!"
}

func (a *Listening) ReadingAnswer() {
	arr := make([]string, 0)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		// Scans a line from Stdin(Console)
		scanner.Scan()
		// Holds the string that scanned
		text := scanner.Text()
		if len(text) != 0 {
			arr = append(arr, text)
			a.AnswerChannel <- respond{
				text: text + " let's do it!",
				ok:   true,
			}
			return
		} else {
			a.AnswerChannel <- respond{}
			return
		}
	}
}

func (a *Listening) waiting(ctx context.Context) respond {
	go a.ReadingAnswer()

	for {
		select {
		case <-ctx.Done():
			return respond{
				text: "",
				ok:   false,
			}
		case rsp := <-a.AnswerChannel:
			return rsp
		}
	}
}

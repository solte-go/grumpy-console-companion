package dictionary

import (
	"grumpy-console-companion/sotle-go/pkg/grcpclient"
	"grumpy-console-companion/sotle-go/pkg/grumpy/dictionary/greetings"
	"grumpy-console-companion/sotle-go/pkg/grumpy/dictionary/standby"
)

type Dictionary struct {
	StandBy   standby.StandBy
	Greetings greetings.Greetings
	client    *grcpclient.Client
}

func New(client *grcpclient.Client) *Dictionary {
	//respond, thought, err := client.QOTD(context.Background(), "greetings")
	//if err != nil {
	//	panic(err)
	//}

	return &Dictionary{
		client: client,
	}
}

//TODO dictionary update by scheduler

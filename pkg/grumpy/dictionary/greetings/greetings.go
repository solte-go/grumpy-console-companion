package greetings

import (
	"grumpy-console-companion/sotle-go/pkg/model"
	"math/rand"
)

type Greetings struct {
	Name   string
	Topics model.Topics
}

func New() *Greetings {
	return &Greetings{}
}

func (g *Greetings) GetRandom() model.Thoughts {
	return g.Topics.Thoughts[rand.Intn(len(g.Topics.Thoughts))]
}

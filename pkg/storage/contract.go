package storage

import (
	"context"
	"grumpy-console-companion/sotle-go/pkg/model"
)

type Storage struct {
	Topics
	Thoughts
}

type Topics interface {
	GetAllTopics(ctx context.Context) ([]model.Topics, error)
}

type Thoughts interface {
	GetThoughtsOnTopic(ctx context.Context, topic string) ([]model.Thoughts, error)
}

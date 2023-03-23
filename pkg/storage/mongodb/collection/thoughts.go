package collection

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"grumpy-console-companion/sotle-go/pkg/model"
)

type ThoughtsCollection struct {
	logger     *zap.Logger
	collection *mongo.Collection
}

// NewThoughtsCollection constructor for thoughts collection
func NewThoughtsCollection(collection *mongo.Collection, logger *zap.Logger) *ThoughtsCollection {
	return &ThoughtsCollection{
		logger:     logger,
		collection: collection,
	}
}

func (t *ThoughtsCollection) GetThoughtsOnTopic(ctx context.Context, topic string) ([]model.Thoughts, error) {
	var result []model.Thoughts
	filter := bson.M{"topic": topic}

	cursor, err := t.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		if err != nil {
			return nil, errors.Wrap(err, "mongodb: error occurred during decoding data from GetThoughtsOnTopic(). error")
		}
	}
	return result, nil
}

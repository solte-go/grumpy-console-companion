package collection

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"grumpy-console-companion/sotle-go/pkg/model"
)

type TopicsCollection struct {
	logger     *zap.Logger
	collection *mongo.Collection
}

// NewTopicsCollection Creates new instance of RequestsCollection.
func NewTopicsCollection(collection *mongo.Collection, logger *zap.Logger) *TopicsCollection {
	return &TopicsCollection{
		logger:     logger,
		collection: collection,
	}
}

func (tc *TopicsCollection) GetAllTopics(ctx context.Context) ([]model.Topics, error) {
	var result []model.Topics
	filter := bson.M{}

	cursor, err := tc.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		if err != nil {
			return nil, errors.Wrap(err, "mongodb: error occurred during decoding data from GetTopics(). error")
		}
	}
	return result, nil
}

func (tc *TopicsCollection) GetTopic(ctx context.Context, topic string) (model.Topics, error) {
	var result model.Topics
	filter := bson.M{"name": topic}

	err := tc.collection.FindOne(ctx, filter).Decode(result)
	if err != nil {
		return model.Topics{}, err
	}

	return result, nil
}

func (tc *TopicsCollection) AddNewTopic(ctx context.Context) ([]string, error) {
	var result []string
	filter := bson.M{}

	cursor, err := tc.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		if err != nil {
			return nil, errors.Wrap(err, "mongodb: error occurred during decoding data from GetTopics(). error")
		}
	}
	return result, nil
}

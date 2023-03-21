package collection

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type ThoughtsCollection struct {
	logger     *zap.Logger
	collection *mongo.Collection
}

// NewThoughtsCollection Creates new instance of RequestsCollection.
func NewThoughtsCollection(collection *mongo.Collection, logger *zap.Logger) *ThoughtsCollection {
	return &ThoughtsCollection{
		logger:     logger,
		collection: collection,
	}
}

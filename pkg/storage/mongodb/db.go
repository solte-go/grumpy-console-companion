package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.uber.org/zap"
	"grumpy-console-companion/sotle-go/config"
	"grumpy-console-companion/sotle-go/pkg/storage/mongodb/collection"
)

type DB struct {
	Thoughts *collection.ThoughtsCollection
}

func InitDatabase(ctx context.Context, conf config.Config, logger *zap.Logger) (*DB, error) {
	db, err := NewDB(ctx, conf.MongoDB)
	if err != nil {
		return nil, err
	}

	instance := &DB{
		Thoughts: collection.NewThoughtsCollection(db.Collection(conf.MongoDB.ThoughtsCollection), logger),
	}

	return instance, nil
}

func NewDB(ctx context.Context, config *config.MongoDB) (*mongo.Database, error) {
	opts := options.Client().
		ApplyURI(config.DatabaseURL).
		SetConnectTimeout(config.ConnectTimeout).
		SetMaxPoolSize(uint64(config.MaxPoolSize)).
		SetWriteConcern(writeconcern.New(writeconcern.WMajority()))

	client, err := mongo.Connect(
		ctx,
		opts,
	)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("ping to mongo failed: %w", err)
	}

	database := client.Database(config.DatabaseName)
	return database, nil
}

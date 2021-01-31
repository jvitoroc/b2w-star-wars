package utils

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DisconnectFunc func() error

func WithTimeout(timeout int64) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
}

func ConnectMongoDB(ctx context.Context, uri string) (*mongo.Client, DisconnectFunc, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	disconnect := func() error { return client.Disconnect(ctx) }
	return client, disconnect, err
}

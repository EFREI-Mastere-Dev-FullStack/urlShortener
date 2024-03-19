package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBConnection struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoDBConnection(uri, dbName string) (*MongoDBConnection, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	database := client.Database(dbName)

	return &MongoDBConnection{
		Client:   client,
		Database: database,
	}, nil
}

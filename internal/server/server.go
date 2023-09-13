package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slog"
)

func connectDB(config *Config) (*mongo.Client, error) {
	connectionString := fmt.Sprintf("mongodb://%s%s", config.DbHost, config.DbPort)

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		return nil, err
	}
	fmt.Printf("Connected to DB on URI: %s\n", connectionString)
	return client, nil
}

func GetDatabase(client *mongo.Client, name string) *mongo.Database {
	return client.Database(name, nil)
}
func GetCollection(client *mongo.Database, name string) *mongo.Collection {
	return client.Collection(name, nil)
}

func StartServer(config *Config) error {
	client, err := connectDB(config)
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())
	server := initServer()
	server.logger.LogAttrs(context.Background(), slog.LevelInfo,
		"Server started",
		slog.Time("at", time.Now()),
		slog.String("port", config.Port),
	)
	return http.ListenAndServe(config.Port, server)
}

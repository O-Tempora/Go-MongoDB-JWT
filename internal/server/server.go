package server

import (
	"context"
	"fmt"
	"gomongojwt/internal/models"
	"gomongojwt/internal/repository"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slog"
)

func connectDB(ctx context.Context, config *Config) (*mongo.Client, error) {
	connectionString := fmt.Sprintf("mongodb://%s%s", config.DbHost, config.DbPort)
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
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

func seedUsers(db *mongo.Database, collectionName string) {
	db.Collection(collectionName).InsertMany(context.Background(), []interface{}{
		models.User{GUID: primitive.NewObjectIDFromTimestamp(time.Now()), Name: "Bonnie"},
		models.User{GUID: primitive.NewObjectIDFromTimestamp(time.Now().Add(2 * time.Minute)), Name: "Clyde"},
	}, nil)
}

func StartServer(config *Config) error {
	ctx := context.Background()
	client, err := connectDB(ctx, config)
	if err != nil {
		return err
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	server := initServer()
	server.client = client
	server.database = server.client.Database(config.Database, nil)
	if err := server.database.Collection(config.Collection).Drop(ctx); err != nil {
		return err
	}
	seedUsers(server.database, config.Collection)
	server.store = repository.CreateStore(server.database)

	server.logger.LogAttrs(ctx, slog.LevelInfo,
		"Server started",
		slog.Time("at", time.Now()),
		slog.String("port", config.Port),
	)

	return http.ListenAndServe(config.Port, server)
}

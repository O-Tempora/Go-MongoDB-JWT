package main

import (
	"context"
	"flag"
	"fmt"
	"gomongojwt/internal/server"
	"io"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v3"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "./configs/default.yaml", "server and db configuration")
}

func main() {
	flag.Parse()
	file, err := os.Open(configPath)
	if err != nil {
		log.Fatal("Selected config file doesn't exist")
	}
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Unable to read config file")
	}
	config := server.NewConfig()
	if err = yaml.Unmarshal(data, config); err != nil {
		log.Fatal("Failed to parse config")
	}
	connectionString := fmt.Sprintf("mongodb://%s%s", config.DbHost, config.DbPort)
	fmt.Printf("connectionString: %v\n", connectionString)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}

	if client.Ping(context.Background(), nil) != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfull connection")

	// collection := client.Database("testing").Collection("numbers")
	// fmt.Println(collection.Name())
	// res, err := collection.InsertOne(context.TODO(), bson.D{{"name", "pi"}, {"value", 3.14159}})
	// fmt.Println(res.InsertedID)

	// var result struct {
	// 	Value float64
	// }
	// cur := collection.FindOne(context.Background(), bson.D{{"name", "pi"}}).Decode(&result)
	// fmt.Println(result, "\n", cur)
}

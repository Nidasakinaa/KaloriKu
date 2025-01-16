package config

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func Init() {
    // Load .env file
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Get MongoDB connection string from .env file
    mongoURI := os.Getenv("MONGOSTRING")

    // Connect to MongoDB
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	// Set global MongoClient
	MongoClient = client
}


func GetMongoClient() *mongo.Client {
	return MongoClient
}
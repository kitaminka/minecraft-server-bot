package database

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var (
	MongoClient       mongo.Client
	MongoDatabaseName string
)

func Connect(mongoUri, mongoDatabaseName string) {
	mongoClient, err := mongo.Connect(nil, options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	log.Print("Successfully connected to MongoDB")

	MongoClient = *mongoClient
	MongoDatabaseName = mongoDatabaseName
}

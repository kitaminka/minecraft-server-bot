package db

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const (
	MongoDatabaseName    = "server-bot"
	MemberCollectionName = "members"
)

var MongoClient mongo.Client

type ServerMember struct {
	ID string
}

func Connect(mongoUri string) {
	mongoClient, err := mongo.Connect(nil, options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	log.Print("Successfully connected to MongoDB")

	MongoClient = *mongoClient
}

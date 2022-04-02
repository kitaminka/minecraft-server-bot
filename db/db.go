package db

import (
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
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
		log.Panicf("Error connecting to MongoDB: %v", err)
	}

	log.Print("Successfully connected to MongoDB")

	MongoClient = *mongoClient
}
func CreateNewMember(member discordgo.Member) {
	_, exists := GetMember(member.User.ID)

	if exists {
		return
	}

	collection := MongoClient.Database(MongoDatabaseName).Collection(MemberCollectionName)

	_, err := collection.InsertOne(nil, bson.D{{"id", member.User.ID}})
	if err != nil {
		log.Printf("Error creating new member: %v", err)
	}
}
func GetMember(id string) (ServerMember, bool) {
	var serverMember ServerMember

	collection := MongoClient.Database(MongoDatabaseName).Collection(MemberCollectionName)

	result := collection.FindOne(nil, bson.D{{"id", id}})

	err := result.Decode(&serverMember)
	if err != nil {
		return ServerMember{}, false
	}

	return serverMember, true
}

package database

import (
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func CreateNewMember(member discordgo.Member) {
	collection := MongoClient.Database(MongoDatabaseName).Collection("members")

	_, err := collection.InsertOne(nil, bson.D{{"id", member.User.ID}})
	if err != nil {
		log.Printf("Error creating new member: %v", err)
	}
}

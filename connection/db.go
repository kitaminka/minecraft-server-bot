package connection

import (
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const (
	Chars                = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	MongoDatabaseName    = "minecraft-server-bot"
	MemberCollectionName = "players"
)

var MongoClient *mongo.Client

type Player struct {
	ID                string
	MinecraftNickname string
}

func ConnectMongo(mongoUri string) {
	mongoClient, err := mongo.Connect(nil, options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Panicf("Error connecting to MongoDB: %v", err)
	}

	log.Print("Successfully connected to MongoDB")

	MongoClient = mongoClient
}
func CreateNewPlayer(member *discordgo.Member, minecraftNickname string) bool {
	_, exists := GetPlayer(member.User.ID)

	if exists {
		return true
	}

	collection := MongoClient.Database(MongoDatabaseName).Collection(MemberCollectionName)

	_, err := collection.InsertOne(nil, bson.D{{"id", member.User.ID}, {"minecraftNickname", minecraftNickname}})
	if err != nil {
		log.Printf("Error creating new member: %v", err)
	}

	return false
}
func GetPlayer(id string) (Player, bool) {
	var serverMember Player

	collection := MongoClient.Database(MongoDatabaseName).Collection(MemberCollectionName)

	result := collection.FindOne(nil, bson.D{{"id", id}})

	err := result.Decode(&serverMember)
	if err != nil {
		return Player{}, false
	}

	return serverMember, true
}

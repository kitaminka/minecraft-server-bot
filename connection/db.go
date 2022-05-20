package connection

import (
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// TODO Move constants to config
// TODO Change all bool to error

const (
	Chars                = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	MongoDatabaseName    = "minecraft-server-bot"
	MemberCollectionName = "players"
)

var MongoClient *mongo.Client

type Player struct {
	DiscordId         string
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
func CreatePlayer(member *discordgo.Member, minecraftNickname string) error {
	_, errDiscord := GetPlayerByDiscord(member)
	_, errMinecraft := GetPlayerByMinecraft(minecraftNickname)

	if errDiscord != nil {
		return errDiscord
	} else if errMinecraft != nil {
		return errMinecraft
	}

	collection := MongoClient.Database(MongoDatabaseName).Collection(MemberCollectionName)

	_, err := collection.InsertOne(nil, bson.D{{"discordId", member.User.ID}, {"minecraftNickname", minecraftNickname}})
	if err != nil {
		log.Printf("Error creating player: %v", err)
		return err
	}

	return nil
}
func DeletePlayer(member *discordgo.Member) error {
	collection := MongoClient.Database(MongoDatabaseName).Collection(MemberCollectionName)

	_, err := collection.DeleteOne(nil, bson.D{{"discordId", member.User.ID}})
	if err != nil {
		log.Printf("Error deleting player: %v", err)
		return err
	}

	return nil
}
func GetPlayerByDiscord(member *discordgo.Member) (Player, error) {
	var serverPlayer Player

	collection := MongoClient.Database(MongoDatabaseName).Collection(MemberCollectionName)

	result := collection.FindOne(nil, bson.D{{"discordId", member.User.ID}})

	err := result.Decode(&serverPlayer)
	if err != nil {
		return Player{}, err
	}

	return serverPlayer, nil
}
func GetPlayerByMinecraft(minecraftNickname string) (Player, error) {
	var serverMember Player

	collection := MongoClient.Database(MongoDatabaseName).Collection(MemberCollectionName)

	result := collection.FindOne(nil, bson.D{{"minecraftNickname", minecraftNickname}})

	err := result.Decode(&serverMember)
	if err != nil {
		return Player{}, err
	}

	return serverMember, nil
}

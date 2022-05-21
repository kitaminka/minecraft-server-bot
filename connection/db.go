package connection

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/minecraft-server-bot/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var MongoClient *mongo.Client

type Player struct {
	DiscordId         string
	MinecraftNickname string
}

func ConnectMongo() {
	mongoClient, err := mongo.Connect(nil, options.Client().ApplyURI(config.Config.Mongo.Uri))
	if err != nil {
		log.Panicf("Error connecting to MongoDB: %v", err)
	}

	log.Print("Successfully connected to MongoDB")

	MongoClient = mongoClient
}
func CreatePlayer(member *discordgo.Member, minecraftNickname string) error {
	_, errDiscord := GetPlayerByDiscord(member)
	_, errMinecraft := GetPlayerByMinecraft(minecraftNickname)

	if errDiscord == nil {
		return fmt.Errorf("player already exists")
	} else if errMinecraft == nil {
		return fmt.Errorf("player already exists")
	}

	collection := MongoClient.Database(config.Config.Mongo.Database).Collection(config.Config.Mongo.PlayerCollection)

	_, err := collection.InsertOne(nil, bson.D{{"discordId", member.User.ID}, {"minecraftNickname", minecraftNickname}})
	if err != nil {
		log.Printf("Error creating player: %v", err)
		return err
	}

	return nil
}
func DeletePlayer(member *discordgo.Member) error {
	collection := MongoClient.Database(config.Config.Mongo.Database).Collection(config.Config.Mongo.PlayerCollection)

	_, err := collection.DeleteOne(nil, bson.D{{"discordId", member.User.ID}})
	if err != nil {
		log.Printf("Error deleting player: %v", err)
		return err
	}

	return nil
}
func GetPlayerByDiscord(member *discordgo.Member) (Player, error) {
	var serverPlayer Player

	collection := MongoClient.Database(config.Config.Mongo.Database).Collection(config.Config.Mongo.PlayerCollection)

	result := collection.FindOne(nil, bson.D{{"discordId", member.User.ID}})

	err := result.Decode(&serverPlayer)
	if err != nil {
		return Player{}, err
	}

	return serverPlayer, nil
}
func GetPlayerByMinecraft(minecraftNickname string) (Player, error) {
	var serverMember Player

	collection := MongoClient.Database(config.Config.Mongo.Database).Collection(config.Config.Mongo.PlayerCollection)

	result := collection.FindOne(nil, bson.D{{"minecraftNickname", minecraftNickname}})

	err := result.Decode(&serverMember)
	if err != nil {
		return Player{}, err
	}

	return serverMember, nil
}

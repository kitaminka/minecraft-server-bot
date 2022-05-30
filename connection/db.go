package connection

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const (
	MongoDatabase           = "minecraft-server-bot"
	MongoPlayerCollection   = "players"
	MongoSettingsCollection = "settings"
)

var MongoClient *mongo.Client

type Player struct {
	DiscordId         string
	MinecraftNickname string
}
type Setting struct {
	Name  string
	Value string
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

	if errDiscord == nil {
		return fmt.Errorf("player already exists")
	} else if errMinecraft == nil {
		return fmt.Errorf("player already exists")
	}

	collection := MongoClient.Database(MongoDatabase).Collection(MongoPlayerCollection)

	_, err := collection.InsertOne(nil, bson.D{{"discordId", member.User.ID}, {"minecraftNickname", minecraftNickname}})
	if err != nil {
		log.Printf("Error creating player: %v", err)
		return err
	}

	return nil
}
func DeletePlayer(member *discordgo.Member) error {
	collection := MongoClient.Database(MongoDatabase).Collection(MongoPlayerCollection)

	_, err := collection.DeleteOne(nil, bson.D{{"discordId", member.User.ID}})
	if err != nil {
		log.Printf("Error deleting player: %v", err)
		return err
	}

	return nil
}
func GetPlayerByDiscord(member *discordgo.Member) (Player, error) {
	var serverPlayer Player

	collection := MongoClient.Database(MongoDatabase).Collection(MongoPlayerCollection)

	result := collection.FindOne(nil, bson.D{{"discordId", member.User.ID}})

	err := result.Decode(&serverPlayer)
	if err != nil {
		return Player{}, err
	}

	return serverPlayer, nil
}
func GetPlayerByMinecraft(minecraftNickname string) (Player, error) {
	var serverMember Player

	collection := MongoClient.Database(MongoDatabase).Collection(MongoPlayerCollection)

	result := collection.FindOne(nil, bson.D{{"minecraftNickname", minecraftNickname}})

	err := result.Decode(&serverMember)
	if err != nil {
		return Player{}, err
	}

	return serverMember, nil
}
func ViewSettings() ([]Setting, error) {
	var settings []Setting

	collection := MongoClient.Database(MongoDatabase).Collection(MongoSettingsCollection)

	result, _ := collection.Find(nil, bson.D{})

	err := result.All(nil, &settings)
	if err != nil {
		return []Setting{}, err
	}

	return settings, nil
}
func GetSetting(settingName string) (Setting, error) {
	var setting Setting

	collection := MongoClient.Database(MongoDatabase).Collection(MongoSettingsCollection)

	result := collection.FindOne(nil, bson.D{{"name", settingName}})

	err := result.Decode(&setting)
	if err != nil {
		return Setting{}, err
	}

	return setting, nil
}
func SetSettingValue(settingName, settingValue string) error {
	collection := MongoClient.Database(MongoDatabase).Collection(MongoSettingsCollection)

	replaceResult, err := collection.ReplaceOne(nil, bson.D{{"name", settingName}}, bson.D{{"name", settingName}, {"value", settingValue}})
	if err != nil {
		return err
	}

	if replaceResult.ModifiedCount == 0 {
		_, err := collection.InsertOne(nil, bson.D{{"name", settingName}, {"value", settingValue}})
		if err != nil {
			return err
		}
	}

	return nil
}
func DeleteSetting(settingName string) error {
	collection := MongoClient.Database(MongoDatabase).Collection(MongoSettingsCollection)

	result, err := collection.DeleteOne(nil, bson.D{{"name", settingName}})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("setting not found")
	}

	return nil
}

package connection

import (
	"errors"
	"fmt"
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

type SettingName string

const (
	MinecraftRoleSetting       SettingName = "minecraftRole"
	WhitelistChannelSetting    SettingName = "whitelistChannel"
	WhitelistMessageSetting    SettingName = "whitelistMessage"
	ApplicationCategorySetting SettingName = "applicationCategory"
	CuratorRoleSetting         SettingName = "curatorRole"
)

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
func GetPlayers() ([]Player, error) {
	var players []Player

	collection := MongoClient.Database(MongoDatabase).Collection(MongoPlayerCollection)
	result, _ := collection.Find(nil, bson.D{})
	err := result.All(nil, &players)

	return players, err
}
func CreatePlayer(userId string, minecraftNickname string) error {
	_, errDiscord := GetPlayerByDiscord(userId)
	_, errMinecraft := GetPlayerByMinecraft(minecraftNickname)

	if errDiscord == nil {
		return fmt.Errorf("player already exists")
	} else if errMinecraft == nil {
		return fmt.Errorf("player already exists")
	}

	collection := MongoClient.Database(MongoDatabase).Collection(MongoPlayerCollection)

	_, err := collection.InsertOne(nil, bson.D{{"discordId", userId}, {"minecraftNickname", minecraftNickname}})
	if err != nil {
		log.Printf("Error creating player: %v", err)
		return err
	}

	return nil
}
func DeletePlayer(userId string) error {
	collection := MongoClient.Database(MongoDatabase).Collection(MongoPlayerCollection)

	_, err := collection.DeleteOne(nil, bson.D{{"discordId", userId}})
	if err != nil {
		log.Printf("Error deleting player: %v", err)
		return err
	}

	return nil
}
func GetPlayerByDiscord(userId string) (Player, error) {
	var serverPlayer Player

	collection := MongoClient.Database(MongoDatabase).Collection(MongoPlayerCollection)
	result := collection.FindOne(nil, bson.D{{"discordId", userId}})
	err := result.Decode(&serverPlayer)

	return serverPlayer, err
}
func GetPlayerByMinecraft(minecraftNickname string) (Player, error) {
	var serverMember Player

	collection := MongoClient.Database(MongoDatabase).Collection(MongoPlayerCollection)
	result := collection.FindOne(nil, bson.D{{"minecraftNickname", minecraftNickname}})
	err := result.Decode(&serverMember)

	return serverMember, err
}
func GetPlayerCount() (int, error) {
	collection := MongoClient.Database(MongoDatabase).Collection(MongoPlayerCollection)
	count, err := collection.CountDocuments(nil, bson.D{})

	return int(count), err
}
func GetSettings() ([]Setting, error) {
	var settings []Setting

	collection := MongoClient.Database(MongoDatabase).Collection(MongoSettingsCollection)

	result, _ := collection.Find(nil, bson.D{})

	err := result.All(nil, &settings)

	return settings, err
}
func GetSetting(settingName SettingName) (Setting, error) {
	var setting Setting

	collection := MongoClient.Database(MongoDatabase).Collection(MongoSettingsCollection)
	result := collection.FindOne(nil, bson.D{{"name", settingName}})
	err := result.Decode(&setting)

	return setting, err
}
func SetSettingValue(settingName SettingName, settingValue string) error {
	collection := MongoClient.Database(MongoDatabase).Collection(MongoSettingsCollection)

	replaceResult, err := collection.ReplaceOne(nil, bson.D{{"name", settingName}}, bson.D{{"name", settingName}, {"value", settingValue}})
	if err != nil {
		return err
	}

	if replaceResult.ModifiedCount == 0 {
		_, err = collection.InsertOne(nil, bson.D{{"name", settingName}, {"value", settingValue}})
		if err != nil {
			return err
		}
	}

	return nil
}
func DeleteSetting(settingName SettingName) error {
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

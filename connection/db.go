package connection

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const (
	MongoDatabase                          = "minecraft-server-bot"
	MongoPlayerCollection                  = "players"
	MongoSettingsCollection                = "settings"
	MinecraftRoleSetting       SettingName = "minecraftRole"
	WhitelistChannelSetting    SettingName = "whitelistChannel"
	WhitelistMessageSetting    SettingName = "whitelistMessage"
	ApplicationCategorySetting SettingName = "applicationCategory"
)

var (
	MongoClient              *mongo.Client
	PlayerCollection         *mongo.Collection
	SettingsCollection       *mongo.Collection
	PlayerAlreadyExistsError = errors.New("player already exists")
	SettingNotFoundError     = errors.New("setting not found")
)

type Player struct {
	DiscordId         string
	MinecraftNickname string
}

type SettingName string

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
	PlayerCollection = MongoClient.Database(MongoDatabase).Collection(MongoPlayerCollection)
	SettingsCollection = MongoClient.Database(MongoDatabase).Collection(MongoSettingsCollection)
}
func CreatePlayer(userId string, minecraftNickname string) error {
	_, errDiscord := GetPlayerByDiscord(userId)
	_, errMinecraft := GetPlayerByMinecraft(minecraftNickname)

	if errDiscord == nil {
		return PlayerAlreadyExistsError
	} else if errMinecraft == nil {
		return PlayerAlreadyExistsError
	}

	_, err := PlayerCollection.InsertOne(nil, bson.D{{"discordId", userId}, {"minecraftNickname", minecraftNickname}})
	if err != nil {
		log.Printf("Error creating player: %v", err)
		return err
	}

	return nil
}
func DeletePlayer(userId string) error {
	_, err := PlayerCollection.DeleteOne(nil, bson.D{{"discordId", userId}})
	if err != nil {
		log.Printf("Error deleting player: %v", err)
		return err
	}

	return nil
}
func GetPlayerByDiscord(userId string) (Player, error) {
	var serverPlayer Player

	result := PlayerCollection.FindOne(nil, bson.D{{"discordId", userId}})
	err := result.Decode(&serverPlayer)

	return serverPlayer, err
}
func GetPlayerByMinecraft(minecraftNickname string) (Player, error) {
	var serverMember Player

	result := PlayerCollection.FindOne(nil, bson.D{{"minecraftNickname", minecraftNickname}})
	err := result.Decode(&serverMember)

	return serverMember, err
}
func GetPlayerCount() (int, error) {
	count, err := PlayerCollection.CountDocuments(nil, bson.D{})

	return int(count), err
}
func GetSettings() ([]Setting, error) {
	var settings []Setting

	result, _ := SettingsCollection.Find(nil, bson.D{})

	err := result.All(nil, &settings)

	return settings, err
}
func GetSetting(settingName SettingName) (Setting, error) {
	var setting Setting

	result := SettingsCollection.FindOne(nil, bson.D{{"name", settingName}})
	err := result.Decode(&setting)

	return setting, err
}
func SetSettingValue(settingName SettingName, settingValue string) error {
	replaceResult, err := SettingsCollection.ReplaceOne(nil, bson.D{{"name", settingName}}, bson.D{{"name", settingName}, {"value", settingValue}})
	if err != nil {
		return err
	}

	if replaceResult.ModifiedCount == 0 {
		_, err = SettingsCollection.InsertOne(nil, bson.D{{"name", settingName}, {"value", settingValue}})
		if err != nil {
			return err
		}
	}

	return nil
}
func DeleteSetting(settingName SettingName) error {
	result, err := SettingsCollection.DeleteOne(nil, bson.D{{"name", settingName}})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return SettingNotFoundError
	}

	return nil
}

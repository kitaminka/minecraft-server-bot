package config

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"os"
)

var Config Configuration

type Configuration struct {
	Token          string
	Guild          string
	RemoveCommands bool             `json:"removeCommands"`
	Intents        discordgo.Intent `json:"intents"`
	Rcon           struct {
		Address  string
		Password string
	}
	EmbedColors struct {
		Primary int `json:"primary"`
		Error   int `json:"error"`
	} `json:"embedColors"`
	Roles struct {
		Admin string `json:"admin"`
	} `json:"roles"`
	Channels struct {
		WhitelistInfo string `json:"whitelistInfo"`
	} `json:"channels"`
	Mongo struct {
		Uri              string
		Database         string `json:"database"`
		MemberCollection string `json:"memberCollection"`
	} `json:"mongo"`
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Panicf("Error loading .env file: %v", err)
	}

	Config.Token = os.Getenv("DISCORD_TOKEN")
	Config.Mongo.Uri = os.Getenv("MONGODB_URI")
	Config.Rcon.Address = os.Getenv("RCON_ADDRESS")
	Config.Rcon.Password = os.Getenv("RCON_PASSWORD")
	Config.Guild = os.Getenv("GUILD_ID")
	Config.Roles.Admin = os.Getenv("ADMIN_ROLE_ID")

	log.Print("Successfully loaded .env file")
}
func LoadConfig() {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		log.Panicf("Error opening config.json file: %v", err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &Config)
	if err != nil {
		log.Panicf("Error loading config: %v", err)
	}

	log.Print("Successfully loaded config")
}

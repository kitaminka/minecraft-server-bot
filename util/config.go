package util

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"os"
)

var Config ConfigType

type ConfigType struct {
	Guild    string           `json:"guild"`
	Intents  discordgo.Intent `json:"intents"`
	Channels struct {
		WelcomeMessageChannel string `json:"welcomeMessageChannel"`
	} `json:"channels"`
}

func LoadConfig() {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("Error opening config.json file: %v", err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &Config)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	if Config.Guild == "" {
		log.Fatalf("Guild is not set in config")
	}

	log.Print("Successfully loaded config")
}

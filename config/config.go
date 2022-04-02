package config

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"os"
)

var Config Configuration

type Configuration struct {
	RemoveApplicationCommands bool             `json:"removeApplicationCommands"`
	Guild                     string           `json:"guild"`
	Intents                   discordgo.Intent `json:"intents"`
	Channels                  struct {
		WelcomeMessageChannel string `json:"welcomeMessageChannel"`
	} `json:"channels"`
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

	if Config.Guild == "" {
		log.Panicf("Guild is not set in config")
	}

	log.Print("Successfully loaded config")
}

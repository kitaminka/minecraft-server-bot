package messages

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"os"
)

type Messages struct {
	WelcomeMessageEmbed discordgo.MessageEmbed `json:"welcome_message_embed"`
}

var DiscordMessages = LoadMessages("messages.json")

func LoadMessages(path string) Messages {
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening messages: %v", err)
	}
	log.Printf("Successfully opened %v", path)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var discordMessages Messages

	json.Unmarshal(byteValue, &discordMessages)

	return discordMessages
}

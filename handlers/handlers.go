package handlers

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

func Ready(ready *discordgo.Ready) {
	log.Printf("Logged in as: %v#%v", ready.User.Username, ready.User.Discriminator)
}

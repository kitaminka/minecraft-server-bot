package handlers

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

func Ready(_ *discordgo.Session, ready *discordgo.Ready) {
	log.Printf("Logged in as: %v#%v", ready.User.Username, ready.User.Discriminator)
}

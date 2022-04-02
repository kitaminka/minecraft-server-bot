package handlers

import (
	"github.com/bwmarrin/discordgo"
)

func InteractionCreate(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	if interactionCreate.Type.String() == "ApplicationCommand" {
		Commands[interactionCreate.ApplicationCommandData().Name].Handler(session, interactionCreate)
	}
}

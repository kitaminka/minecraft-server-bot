package handlers

import (
	"github.com/bwmarrin/discordgo"
)

func InteractionCreate(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	if interactionCreate.Type.String() == "ApplicationCommand" {
		Commands[interactionCreate.ID].Handler(session, interactionCreate)
	}
}

package handlers

import (
	"github.com/bwmarrin/discordgo"
)

var Handlers = []interface{}{
	func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
		if interactionCreate.Type.String() == "ApplicationCommand" {
			Commands[interactionCreate.ApplicationCommandData().Name].Handler(session, interactionCreate)
		}
	},
}

func AddHandlers(session *discordgo.Session) {
	for _, value := range Handlers {
		session.AddHandler(value)
	}
}

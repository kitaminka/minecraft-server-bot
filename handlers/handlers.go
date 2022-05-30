package handlers

import (
	"github.com/bwmarrin/discordgo"
)

var GuildId string

const (
	PrimaryEmbedColor = 9383347
	ErrorEmbedColor   = 13179932
)

var Handlers = []interface{}{
	func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
		if interactionCreate.Type == discordgo.InteractionApplicationCommand {
			Commands[interactionCreate.ApplicationCommandData().Name].Handler(session, interactionCreate)
		} else if interactionCreate.Type == discordgo.InteractionMessageComponent {
			Components[interactionCreate.MessageComponentData().CustomID].Handler(session, interactionCreate)
		}
	},
	func(session *discordgo.Session, ready *discordgo.Ready) {
		updateWhitelistMessage(session)
	},
}

func AddHandlers(session *discordgo.Session) {
	for _, value := range Handlers {
		session.AddHandler(value)
	}
}

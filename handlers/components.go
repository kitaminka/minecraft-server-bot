package handlers

import (
	"github.com/bwmarrin/discordgo"
)

type Component struct {
	MessageComponent discordgo.MessageComponent
	Handler          func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate)
}

var Components = map[string]Component{
	"reset-password": {
		MessageComponent: &discordgo.Button{
			CustomID: "reset-password",
			Label:    "Reset password",
			Style:    discordgo.PrimaryButton,
			Emoji: discordgo.ComponentEmoji{
				Name: "🔐",
			},
		},
		Handler: func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
			resetPasswordHandler(session, interactionCreate)
		},
	},
}

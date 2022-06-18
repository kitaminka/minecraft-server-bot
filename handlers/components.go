package handlers

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

type Component struct {
	MessageComponent discordgo.MessageComponent
	Handler          func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate)
}

var Components = map[string]Component{
	"reset_password": {
		MessageComponent: &discordgo.Button{
			CustomID: "reset_password",
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
	"change_password": {
		MessageComponent: &discordgo.Button{
			CustomID: "change_password",
			Label:    "Change password",
			Style:    discordgo.PrimaryButton,
			Emoji: discordgo.ComponentEmoji{
				Name: "✏",
			},
		},
		Handler: func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
			err := session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseModal,
				Data: Modals["change_password"].Modal,
			})
			if err != nil {
				log.Printf("Error responding to interaction: %v", err)
			}
		},
	},
}

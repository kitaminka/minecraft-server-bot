package handlers

import (
	"fmt"
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
				Data: &discordgo.InteractionResponseData{
					Title:    "Change password",
					CustomID: fmt.Sprintf("change_password_%v", interactionCreate.Member.User.ID),
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:    "new_password",
									Label:       "Enter your new password",
									Style:       discordgo.TextInputShort,
									Placeholder: "New password",
									Required:    true,
									MaxLength:   15,
									MinLength:   3,
								},
							},
						},
					},
				},
			})
			if err != nil {
				log.Printf("Error responding to interaction: %v", err)
			}
		},
	},
}

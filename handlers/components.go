package handlers

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/minecraft-server-bot/connection"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type Component struct {
	MessageComponent discordgo.MessageComponent
	Handler          func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate)
}

var Components = map[string]Component{
	"apply_for_whitelist": {
		MessageComponent: &discordgo.Button{
			CustomID: "apply_for_whitelist",
			Label:    "Apply for whitelist",
			Style:    discordgo.SuccessButton,
			Emoji: discordgo.ComponentEmoji{
				Name: "✅",
			},
		},
		Handler: func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
			_, err := connection.GetPlayerByDiscord(interactionCreate.Member.User.ID)
			if !errors.Is(err, mongo.ErrNoDocuments) {
				interactionRespondError(session, interactionCreate.Interaction, "You are already registered.")
				return
			}
			if err != nil {
				log.Printf("Error occurred getting player: %v", err)
				return
			}

			err = session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseModal,
				Data: Modals["apply_for_whitelist"].Modal,
			})
			if err != nil {
				log.Printf("Error responding to interaction: %v", err)
			}
		},
	},
	"reset_password": {
		MessageComponent: &discordgo.Button{
			CustomID: "reset_password",
			Label:    "Reset password",
			Style:    discordgo.PrimaryButton,
			Emoji: discordgo.ComponentEmoji{
				Name: "🔐",
			},
		},
		Handler: resetPasswordHandler,
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

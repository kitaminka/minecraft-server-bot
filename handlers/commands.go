package handlers

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

type Command struct {
	ApplicationCommand discordgo.ApplicationCommand
	Handler            func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate)
}

var Commands = map[string]Command{
	"ping": {
		ApplicationCommand: discordgo.ApplicationCommand{
			Name:        "ping",
			Description: "Тестовая команда для проверки работы бота",
		},
		Handler: func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
			err := session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   1 << 6,
					Content: "Pong!",
				},
			})

			if err != nil {
				log.Printf("Error responding to interaction: %v", err)
			}
		},
	},
}

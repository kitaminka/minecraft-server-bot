package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/server-bot/config"
	"log"
)

type Command struct {
	ApplicationCommand *discordgo.ApplicationCommand
	Handler            func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate)
}

var Commands = map[string]Command{
	"ping": {
		ApplicationCommand: &discordgo.ApplicationCommand{
			Type:        discordgo.ChatApplicationCommand,
			Name:        "ping",
			Description: "Тестовая команда",
		},
		Handler: func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
			err := session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Pong!",
					Flags:   1 << 6,
				},
			})

			if err != nil {
				log.Printf("Error responding to interaction: %v", err)
			}
		},
	},
}

func CreateApplicationCommands(session *discordgo.Session) {
	for index, value := range Commands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, config.Config.Guild, value.ApplicationCommand)
		if err != nil {
			log.Panicf("Error creating '%v' command: %v", value.ApplicationCommand.Name, err)
		}
		log.Printf("Successfully created '%v' command", cmd.Name)

		if command, exists := Commands[index]; exists {
			command.ApplicationCommand = cmd
			Commands[index] = command
		}
	}
}

func RemoveApplicationCommands(session *discordgo.Session) {
	for _, value := range Commands {
		err := session.ApplicationCommandDelete(session.State.User.ID, config.Config.Guild, value.ApplicationCommand.ID)
		if err != nil {
			log.Panicf("Error deleting '%v' command: %v", value.ApplicationCommand.Name, err)
		}
		log.Printf("Successfully deleted '%v' command", value.ApplicationCommand.Name)
	}
}

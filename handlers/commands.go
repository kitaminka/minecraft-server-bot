package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/minecraft-server-bot/config"
	"github.com/kitaminka/minecraft-server-bot/connection"
	"log"
	"strings"
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
			Description: "Pong!",
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
	"whitelist": {
		ApplicationCommand: &discordgo.ApplicationCommand{
			Type:        discordgo.ChatApplicationCommand,
			Name:        "whitelist",
			Description: "Get player whitelist",
		},
		Handler: func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
			players := connection.GetPlayerWhitelist()
			err := session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Whitelist",
							Description: "Player nicknames: **" + strings.Join(players, ", ") + "**",
							Timestamp:   "",
							Color:       config.Config.EmbedColors.Primary,
						},
					},
				},
			})

			if err != nil {
				log.Printf("Error responding to interaction: %v", err)
			}
		},
	},
	"register": {
		ApplicationCommand: &discordgo.ApplicationCommand{
			Type:        discordgo.ChatApplicationCommand,
			Name:        "register",
			Description: "Register new player",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "nickname",
					Description: "Minecraft nickname",
				},
			},
		},
		Handler: func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
			if interactionCreate.Member.Permissions&discordgo.PermissionAdministrator == 0 {
				interactionRespondError(session, interactionCreate.Interaction, "Sorry, you don't have permission.")
				return
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

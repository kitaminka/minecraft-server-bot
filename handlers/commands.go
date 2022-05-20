package handlers

import (
	"fmt"
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
			players, err := connection.GetPlayerWhitelist()
			if err != nil {
				interactionRespondError(session, interactionCreate.Interaction, "Error occurred getting whitelist.")
			}

			err = session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Whitelist",
							Description: fmt.Sprintf("Player nicknames: **%v**", strings.Join(players, ", ")),
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
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "member",
					Description: "Discord server member",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "nickname",
					Description: "Minecraft nickname",
					Required:    true,
				},
			},
		},
		Handler: func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
			if interactionCreate.Member.Permissions&discordgo.PermissionAdministrator == 0 {
				interactionRespondError(session, interactionCreate.Interaction, "Sorry, you don't have permission.")
				return
			}

			options := interactionCreate.ApplicationCommandData().Options
			minecraftNickname := options[1].StringValue()
			member, err := session.GuildMember(config.Config.Guild, options[0].UserValue(session).ID)
			if err != nil {
				log.Printf("Error getting member %v", err)
				return
			}

			err = connection.CreatePlayer(member, minecraftNickname)
			if err != nil {
				interactionRespondError(session, interactionCreate.Interaction, "Error occurred creating player.")
				return
			}

			err = connection.AddPlayerWhitelist(minecraftNickname)
			if err != nil {
				interactionRespondError(session, interactionCreate.Interaction, "Error occurred whitelisting player.")
				connection.DeletePlayer(member)
				return
			}

			password, err := connection.RegisterPlayer(minecraftNickname)
			if err != nil {
				interactionRespondError(session, interactionCreate.Interaction, "Error occurred registering player.")
				connection.DeletePlayer(member)
				connection.RemovePlayerWhitelist(minecraftNickname)
				return
			}

			channel, err := session.UserChannelCreate(member.User.ID)

			if err == nil {
				_, err = session.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Minecraft Server Night Pix",
							Description: "You have been successfully registered on the server.",
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:   "Discord member",
									Value:  fmt.Sprintf("<@%v>", member.User.ID),
									Inline: true,
								},
								{
									Name:   "Minecraft nickname",
									Value:  minecraftNickname,
									Inline: true,
								},
								{
									Name:   "Password",
									Value:  fmt.Sprintf("||%v||", password),
									Inline: true,
								},
							},
							Color: config.Config.EmbedColors.Primary,
						},
					},
				})
			}

			if err != nil {
				log.Printf("Error sending message: %v", err)

				err = session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							{
								Title:       "Player registered",
								Description: "Successfully registered new player.\n**Error occurred sending password!**",
								Fields: []*discordgo.MessageEmbedField{
									{
										Name:   "Discord member",
										Value:  fmt.Sprintf("<@%v>", member.User.ID),
										Inline: true,
									},
									{
										Name:   "Minecraft nickname",
										Value:  minecraftNickname,
										Inline: true,
									},
									{
										Name:   "Password",
										Value:  fmt.Sprintf("||%v||", password),
										Inline: true,
									},
								},
								Color: config.Config.EmbedColors.Primary,
							},
						},
						Flags: 1 << 6,
					},
				})
				if err != nil {
					log.Printf("Error responding to interaction: %v", err)
				}
				return
			}

			err = session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Player registered",
							Description: "Successfully registered new player.",
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:   "Discord member",
									Value:  fmt.Sprintf("<@%v>", member.User.ID),
									Inline: true,
								},
								{
									Name:   "Minecraft nickname",
									Value:  minecraftNickname,
									Inline: true,
								},
								{
									Name:   "Password",
									Value:  fmt.Sprintf("||%v||", password),
									Inline: true,
								},
							},
							Color: config.Config.EmbedColors.Primary,
						},
					},
					Flags: 1 << 6,
				},
			})
			if err != nil {
				log.Printf("Error responding to interaction: %v", err)
			}
		},
	},
	"unregister": {
		ApplicationCommand: &discordgo.ApplicationCommand{
			Type:        discordgo.ChatApplicationCommand,
			Name:        "unregister",
			Description: "Unregister player",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "member",
					Description: "Discord server member",
					Required:    true,
				},
			},
		},
		Handler: func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
			if interactionCreate.Member.Permissions&discordgo.PermissionAdministrator == 0 {
				interactionRespondError(session, interactionCreate.Interaction, "Sorry, you don't have permission.")
				return
			}

			// TODO Unregister player
			// TODO Send reply to interaction
		},
	},
	// TODO Add reset-password command
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

package handlers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
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
			whitelistPlayers, err := connection.GetPlayerWhitelist()
			if err != nil {
				interactionRespondError(session, interactionCreate.Interaction, "Failed to get player whitelist")
				return
			}
			whitelistPlayerCount := len(whitelistPlayers)

			var whitelistPlayerString string
			if whitelistPlayerCount != 0 {
				whitelistPlayerString = "`" + strings.Join(whitelistPlayers, "`, `") + "`"
				if len(whitelistPlayerString) > 1019 {
					whitelistPlayerString = "`" + strings.Join(whitelistPlayers, "`, `")[:1021] + "`..."
				}
			} else {
				whitelistPlayerString = "The whitelist is empty."
			}

			err = session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Whitelist",
							Description: whitelistPlayerString,
							Color:       PrimaryEmbedColor,
						},
					},
				},
			})

			if err != nil {
				log.Printf("Error responding to interaction: %v", err)
			}
		},
	},
	"settings": {
		ApplicationCommand: &discordgo.ApplicationCommand{
			Type:        discordgo.ChatApplicationCommand,
			Name:        "settings",
			Description: "Settings management",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "set",
					Description: "Set setting value",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "name",
							Description: "Setting name",
							Required:    true,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{
									Name:  "Minecraft role ID",
									Value: connection.MinecraftRoleSetting,
								},
								{
									Name:  "Whitelist info channel ID",
									Value: connection.WhitelistChannelSetting,
								},
								{
									Name:  "Whitelist info message ID",
									Value: connection.WhitelistMessageSetting,
								},
							},
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "value",
							Description: "Setting value",
							Required:    true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "view",
					Description: "View all settings",
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "delete",
					Description: "Delete setting",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "name",
							Description: "Setting name",
							Required:    true,
						},
					},
				},
			},
		},
		Handler: func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
			if interactionCreate.Member.Permissions&discordgo.PermissionAdministrator == 0 {
				interactionRespondError(session, interactionCreate.Interaction, "Sorry, you don't have permission.")
				return
			}

			options := interactionCreate.ApplicationCommandData().Options

			if options[0].Name == "set" {
				subcommandOptions := options[0].Options

				settingName := subcommandOptions[0].StringValue()
				settingValue := subcommandOptions[1].StringValue()

				err := connection.SetSettingValue(connection.SettingName(settingName), settingValue)
				if err != nil {
					interactionRespondError(session, interactionCreate.Interaction, fmt.Sprintf("Error occurred: %v", err))
					return
				}

				err = session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							{
								Title:       "Setting edited",
								Description: "Setting edited successfully",
								Color:       PrimaryEmbedColor,
								Fields: []*discordgo.MessageEmbedField{
									{
										Name:  "Name",
										Value: settingName,
									},
									{
										Name:  "Value",
										Value: settingValue,
									},
								},
							},
						},
						Flags: 1 << 6,
					},
				})
				if err != nil {
					log.Printf("Error responding to interaction: %v", err)
				}
			} else if options[0].Name == "view" {
				settings, err := connection.GetSettings()
				if err != nil {
					interactionRespondError(session, interactionCreate.Interaction, fmt.Sprintf("Error occurred: %v", err))
					return
				}

				var fields []*discordgo.MessageEmbedField
				for _, setting := range settings {
					fields = append(fields, &discordgo.MessageEmbedField{
						Name:  setting.Name,
						Value: setting.Value,
					})
				}

				err = session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							{
								Title:       "Settings",
								Description: "Full list of settings",
								Color:       PrimaryEmbedColor,
								Fields:      fields,
							},
						},
						Flags: 1 << 6,
					},
				})
				if err != nil {
					log.Printf("Error responding to interaction: %v", err)
				}
			} else if options[0].Name == "delete" {
				subcommandOptions := options[0].Options
				settingName := subcommandOptions[0].StringValue()

				err := connection.DeleteSetting(connection.SettingName(settingName))
				if err != nil {
					interactionRespondError(session, interactionCreate.Interaction, fmt.Sprintf("Error occurred: %v", err))
					return
				}

				err = session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							{
								Title:       "Setting deleted",
								Description: "Setting deleted successfully",
								Color:       PrimaryEmbedColor,
								Fields: []*discordgo.MessageEmbedField{
									{
										Name:  "Name",
										Value: settingName,
									},
								},
							},
						},
						Flags: 1 << 6,
					},
				})
				if err != nil {
					log.Printf("Error responding to interaction: %v", err)
				}
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

			err := session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags: 1 << 6,
				},
			})
			if err != nil {
				log.Printf("Error responding to interaction: %v", err)
				return
			}

			options := interactionCreate.ApplicationCommandData().Options
			minecraftNickname := options[1].StringValue()
			user := options[0].UserValue(session)

			password, err := connection.RegisterPlayer(options[0].UserValue(session).ID, minecraftNickname)
			if err == connection.PlayerAlreadyExistsError {
				followupErrorMessageCreate(session, interactionCreate.Interaction, "Player already exists")
				return
			} else if err != nil {
				followupErrorMessageCreate(session, interactionCreate.Interaction, fmt.Sprintf("Error occurred registring player: %v", err))
				connection.DeletePlayer(user.ID)
				connection.RemovePlayerWhitelist(minecraftNickname)
				return
			}

			channel, messageErr := session.UserChannelCreate(user.ID)

			if messageErr == nil {
				_, messageErr = session.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Minecraft Server Night Pix",
							Description: "You have been successfully registered on the server.",
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:   "Discord member",
									Value:  fmt.Sprintf("<@%v>", user.ID),
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
							Color: PrimaryEmbedColor,
						},
					},
				})
			}

			setting, roleErr := connection.GetSetting(connection.MinecraftRoleSetting)
			if roleErr == nil {
				roleErr = session.GuildMemberRoleAdd(GuildId, user.ID, setting.Value)
			}

			go updateWhitelistMessage(session)

			_, err = session.FollowupMessageCreate(interactionCreate.Interaction, true, &discordgo.WebhookParams{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Player registered",
						Description: "Successfully registered new player.",
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:   "Discord member",
								Value:  fmt.Sprintf("<@%v>", user.ID),
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
							{
								Name:   "Message error",
								Value:  fmt.Sprint(messageErr),
								Inline: true,
							},
							{
								Name:   "Role error",
								Value:  fmt.Sprint(roleErr),
								Inline: true,
							},
						},
						Color: PrimaryEmbedColor,
					},
				},
				Flags: 1 << 6,
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

			err := session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags: 1 << 6,
				},
			})
			if err != nil {
				log.Printf("Error responding to interaction: %v", err)
				return
			}

			user := interactionCreate.Interaction.ApplicationCommandData().Options[0].UserValue(session)

			player, err := connection.UnregisterPlayer(user.ID)
			if err != nil {
				log.Printf("Error unregistring player: %v", err)
				followupErrorMessageCreate(session, interactionCreate.Interaction, fmt.Sprintf("Error occurred unregistring player: %v", err))
				return
			}

			setting, roleErr := connection.GetSetting(connection.MinecraftRoleSetting)
			if roleErr == nil {
				roleErr = session.GuildMemberRoleRemove(GuildId, user.ID, setting.Value)
			}

			go updateWhitelistMessage(session)

			_, err = session.FollowupMessageCreate(interactionCreate.Interaction, true, &discordgo.WebhookParams{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Player unregistered",
						Description: "Successfully unregistered player.",
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:   "Discord member",
								Value:  fmt.Sprintf("<@%v>", user.ID),
								Inline: true,
							},
							{
								Name:   "Minecraft nickname",
								Value:  player.MinecraftNickname,
								Inline: true,
							},
							{
								Name:   "Role error",
								Value:  fmt.Sprint(roleErr),
								Inline: true,
							},
						},
						Color: PrimaryEmbedColor,
					},
				},
				Flags: 1 << 6,
			})
			if err != nil {
				log.Printf("Error responding to interaction: %v", err)
			}
		},
	},
	"reset-password": {
		ApplicationCommand: &discordgo.ApplicationCommand{
			Type:        discordgo.ChatApplicationCommand,
			Name:        "reset-password",
			Description: "Reset player password",
		},
		Handler: resetPasswordHandler,
	},
	"send-whitelist": {
		ApplicationCommand: &discordgo.ApplicationCommand{
			Type:        discordgo.ChatApplicationCommand,
			Name:        "send-whitelist",
			Description: "Send whitelist message",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "Whitelist channel",
					ChannelTypes: []discordgo.ChannelType{
						discordgo.ChannelTypeGuildText,
					},
					Required: true,
				},
			},
		},
		Handler: func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
			if interactionCreate.Member.Permissions&discordgo.PermissionAdministrator == 0 {
				interactionRespondError(session, interactionCreate.Interaction, "Sorry, you don't have permission.")
				return
			}

			channel := interactionCreate.ApplicationCommandData().Options[0].ChannelValue(session)

			if channel.Type != discordgo.ChannelTypeGuildText {
				interactionRespondError(session, interactionCreate.Interaction, "Wrong channel type.")
				return
			}

			embed, err := createWhitelistEmbed()
			if err != nil {
				log.Printf("Error creating whitelist message: %v", err)
				return
			}

			message, err := session.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
				Embeds: []*discordgo.MessageEmbed{
					embed,
				},
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							Components["reset_password"].MessageComponent,
							Components["change_password"].MessageComponent,
						},
					},
				},
			})
			if err != nil {
				interactionRespondError(session, interactionCreate.Interaction, fmt.Sprintf("Error occured sending whitelist: %v", err))
				return
			}

			err = connection.SetSettingValue(connection.WhitelistChannelSetting, channel.ID)
			if err != nil {
				interactionRespondError(session, interactionCreate.Interaction, fmt.Sprintf("Error occured: %v", err))
				return
			}
			err = connection.SetSettingValue(connection.WhitelistMessageSetting, message.ID)
			if err != nil {
				err := connection.DeleteSetting(connection.WhitelistChannelSetting)
				if err != nil {
					log.Printf("Error deleting setting: %v", err)
				}
				interactionRespondError(session, interactionCreate.Interaction, fmt.Sprintf("Error occured: %v", err))
				return
			}

			err = session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Whitelist info",
							Description: "Message sent.",
							Color:       PrimaryEmbedColor,
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:  "Whitelist channel",
									Value: channel.ID,
								},
								{
									Name:  "Whitelist message",
									Value: message.ID,
								},
							},
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
	"playtime": {
		ApplicationCommand: &discordgo.ApplicationCommand{
			Type:        discordgo.ChatApplicationCommand,
			Name:        "playtime",
			Description: "Get player playtime",
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
			user := interactionCreate.Interaction.ApplicationCommandData().Options[0].UserValue(session)
			player, err := connection.GetPlayerByDiscord(user.ID)
			if err != nil {
				interactionRespondError(session, interactionCreate.Interaction, fmt.Sprintf("Error occurred getting player: %v", err))
				return
			}
			playtime, err := connection.GetPlayerPlaytime(player.MinecraftNickname)
			if err != nil {
				interactionRespondError(session, interactionCreate.Interaction, fmt.Sprintf("Error occurred getting playtime: %v", err))
				return
			}
			err = session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: "Playtime",
							Color: PrimaryEmbedColor,
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:  "Discord member",
									Value: fmt.Sprintf("<@%v>", user.ID),
								},
								{
									Name:  "Minecraft nickname",
									Value: player.MinecraftNickname,
								},
								{
									Name:  "Playtime",
									Value: playtime,
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

func CreateApplicationCommands(session *discordgo.Session, guildId string) {
	GuildId = guildId

	for index, value := range Commands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, GuildId, value.ApplicationCommand)
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
		err := session.ApplicationCommandDelete(session.State.User.ID, GuildId, value.ApplicationCommand.ID)
		if err != nil {
			log.Panicf("Error deleting '%v' command: %v", value.ApplicationCommand.Name, err)
		}
		log.Printf("Successfully deleted '%v' command", value.ApplicationCommand.Name)
	}
}

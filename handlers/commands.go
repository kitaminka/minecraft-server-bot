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
			players, err := connection.GetPlayerWhitelist()
			if err != nil {
				interactionRespondError(session, interactionCreate.Interaction, fmt.Sprintf("Error occurred getting whitelist: %v", err))
				return
			}

			err = session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Whitelist",
							Description: fmt.Sprintf("Player nicknames: **%v**", strings.Join(players, ", ")),
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
									Value: "minecraftRole",
								},
								{
									Name:  "Whitelist info channel ID",
									Value: "whitelistChannel",
								},
								{
									Name:  "Whitelist info message ID",
									Value: "whitelistMessage",
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

				err := connection.SetSettingValue(settingName, settingValue)
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
				settings, err := connection.ViewSettings()
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

				err := connection.DeleteSetting(settingName)
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
			member, err := session.GuildMember(GuildId, options[0].UserValue(session).ID)
			if err != nil {
				log.Printf("Error getting member %v", err)
				followupErrorMessageCreate(session, interactionCreate.Interaction, fmt.Sprintf("Error occurred getting member: %v", err))
				return
			}

			err = connection.CreatePlayer(member, minecraftNickname)
			if err != nil {
				followupErrorMessageCreate(session, interactionCreate.Interaction, fmt.Sprintf("Error occurred creating player: %v", err))
				return
			}

			err = connection.AddPlayerWhitelist(minecraftNickname)
			if err != nil {
				followupErrorMessageCreate(session, interactionCreate.Interaction, fmt.Sprintf("Error occurred whitelisting player: %v", err))
				err = connection.DeletePlayer(member)
				if err != nil {
					log.Printf("Error deleting player: %v", err)
				}
				return
			}

			password, err := connection.RegisterPlayer(minecraftNickname)
			if err != nil {
				followupErrorMessageCreate(session, interactionCreate.Interaction, fmt.Sprintf("Error occurred registering player: %v", err))
				err = connection.DeletePlayer(member)
				if err != nil {
					log.Printf("Error deleting player: %v", err)
				}
				err = connection.RemovePlayerWhitelist(minecraftNickname)
				if err != nil {
					log.Printf("Error removing player from whitelist: %v", err)
				}
				return
			}

			channel, messageErr := session.UserChannelCreate(member.User.ID)

			if messageErr == nil {
				_, messageErr = session.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
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
							Color: PrimaryEmbedColor,
						},
					},
				})
			}

			setting, roleErr := connection.GetSetting("minecraftRole")
			if roleErr == nil {
				roleErr = session.GuildMemberRoleAdd(GuildId, member.User.ID, setting.Value)
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
			member, err := session.GuildMember(GuildId, user.ID)
			if err != nil {
				log.Printf("Error getting member %v", err)
				followupErrorMessageCreate(session, interactionCreate.Interaction, fmt.Sprintf("Error occurred getting member: %v", err))
				return
			}

			player, err := connection.GetPlayerByDiscord(member)
			if err != nil {
				log.Printf("Error getting player: %v", err)
				followupErrorMessageCreate(session, interactionCreate.Interaction, fmt.Sprintf("Error occurred getting player: %v", err))
				return
			}

			unregisterErr := connection.UnregisterPlayer(player.MinecraftNickname)
			whitelistErr := connection.RemovePlayerWhitelist(player.MinecraftNickname)
			playerErr := connection.DeletePlayer(member)
			setting, roleErr := connection.GetSetting("minecraftRole")
			if roleErr == nil {
				roleErr = session.GuildMemberRoleRemove(GuildId, member.User.ID, setting.Value)
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
								Value:  fmt.Sprintf("<@%v>", member.User.ID),
								Inline: true,
							},
							{
								Name:   "Minecraft nickname",
								Value:  player.MinecraftNickname,
								Inline: true,
							},
							{
								Name:   "Unregister error",
								Value:  fmt.Sprint(unregisterErr),
								Inline: true,
							},
							{
								Name:   "Whitelist error",
								Value:  fmt.Sprint(whitelistErr),
								Inline: true,
							},
							{
								Name:   "Player error",
								Value:  fmt.Sprint(playerErr),
								Inline: true,
							},
							{
								Name:   "Role error",
								Value:  fmt.Sprint(playerErr),
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
		Handler: func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
			resetPasswordHandler(session, interactionCreate)
		},
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

			players, err := connection.GetPlayers()
			if err != nil {
				interactionRespondError(session, interactionCreate.Interaction, fmt.Sprintf("Error occured getting players: %v", err))
				return
			}

			var fields []*discordgo.MessageEmbedField
			for _, player := range players {
				fields = append(fields, &discordgo.MessageEmbedField{
					Name:  player.MinecraftNickname,
					Value: fmt.Sprintf("<@%v>", player.DiscordId),
				})
			}

			message, err := session.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Whitelist info",
						Description: "All Minecraft Server Night Pix players",
						Color:       PrimaryEmbedColor,
						Fields:      fields,
					},
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

			err = connection.SetSettingValue("whitelistChannel", channel.ID)
			if err != nil {
				interactionRespondError(session, interactionCreate.Interaction, fmt.Sprintf("Error occured: %v", err))
				return
			}
			err = connection.SetSettingValue("whitelistMessage", message.ID)
			if err != nil {
				err := connection.DeleteSetting("whitelistChannel")
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

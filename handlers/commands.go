package handlers

import (
	"fmt"
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
	"send-profile-message": {
		ApplicationCommand: &discordgo.ApplicationCommand{
			Type:        discordgo.ChatApplicationCommand,
			Name:        "send-profile-message",
			Description: "Отправить сообщение для создания профиля",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "channel",
					Description: "Канал, в который необходимо отправить сообщение",
					Type:        discordgo.ApplicationCommandOptionChannel,
					Required:    false,
				},
			},
		},
		Handler: func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
			var channelID string

			if interactionCreate.ApplicationCommandData().Options == nil {
				channelID = interactionCreate.ChannelID
			} else if interactionCreate.ApplicationCommandData().Options[0].ChannelValue(session).Type == 1 {
				channelID = interactionCreate.ApplicationCommandData().Options[0].ChannelValue(session).ID
			} else {
				interactionRespondError(session, interactionCreate.Interaction, "Неправильный тип канала. Выберите текстовый канал.")
				return
			}

			_, err := session.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Создать профиль",
						Description: "Нажмите кнопку ниже, чтобы создать и настроить профиль!",
						Color:       8523465,
					},
				},
				Components: []discordgo.MessageComponent{
					&discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Label:    "Создать профиль",
								Style:    discordgo.PrimaryButton,
								Disabled: false,
								Emoji: discordgo.ComponentEmoji{
									Name: "🔑",
								},
								CustomID: "create-profile",
							},
						},
					},
				},
			})

			if err != nil {
				log.Printf("Error sending profile message: %v", err)
				interactionRespondError(session, interactionCreate.Interaction, fmt.Sprintf("Произошла ошибка при отправке сообщения: **%v**", err))
				return
			}

			err = session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Сообщение успешно отправлено",
							Description: fmt.Sprintf("Сообщение успешно отправлено в канал <#%v>.", channelID),
							Color:       8523465,
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

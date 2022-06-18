package handlers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/minecraft-server-bot/connection"
	"log"
)

type Modal struct {
	Modal   *discordgo.InteractionResponseData
	Handler func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate)
}

var Modals = map[string]Modal{
	"change_password": {
		Modal: &discordgo.InteractionResponseData{
			Title:    "Change password",
			CustomID: "change_password",
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
		Handler: func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
			data := interactionCreate.ModalSubmitData()
			member := interactionCreate.Member

			player, err := connection.GetPlayerByDiscord(member)
			if err != nil {
				interactionRespondError(session, interactionCreate.Interaction, fmt.Sprintf("Error occurred getting player: %v", err))
			}

			password := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
			err = connection.ChangePlayerPassword(player.MinecraftNickname, password)
			if err != nil {
				interactionRespondError(session, interactionCreate.Interaction, fmt.Sprintf("Error occurred changing player password: %v", err))
			}

			channel, messageErr := session.UserChannelCreate(member.User.ID)

			if messageErr == nil {
				_, messageErr = session.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Minecraft Server Night Pix",
							Description: "Your password has been changed.",
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
			if messageErr != nil {
				log.Printf("Error sending message: %v", err)
			}

			err = session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Password changed",
							Description: "Your password has been changed.",
							Color:       PrimaryEmbedColor,
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
									Name:   "Password",
									Value:  fmt.Sprintf("||%v||", password),
									Inline: true,
								},
								{
									Name:   "Message error",
									Value:  fmt.Sprint(messageErr),
									Inline: true,
								},
							},
						},
					},
					Flags: 1 << 6,
				},
			})
			if err != nil {
				log.Printf("Error responding to interaction: %v", err)
				return
			}
		},
	},
}
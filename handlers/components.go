package handlers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/minecraft-server-bot/connection"
	"log"
)

type Component struct {
	MessageComponent discordgo.MessageComponent
	Handler          func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate)
}

var Components = map[string]Component{
	"reset-password": {
		MessageComponent: &discordgo.Button{
			CustomID: "reset-password",
			Label:    "Reset password",
			Style:    discordgo.PrimaryButton,
			Emoji: discordgo.ComponentEmoji{
				Name: "🔐",
			},
		},
		Handler: func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
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

			member := interactionCreate.Member

			player, err := connection.GetPlayerByDiscord(member)
			if err != nil {
				log.Printf("Error getting player: %v", err)
				_, err := followupErrorMessageCreate(session, interactionCreate.Interaction, fmt.Sprintf("Error occurred getting player: %v", err))
				if err != nil {
					log.Printf("Error sending message: %v", err)
				}
				return
			}

			password, err := connection.ResetPlayerPassword(player.MinecraftNickname)
			if err != nil {
				log.Printf("Error resetting player password: %v", err)
				_, err := followupErrorMessageCreate(session, interactionCreate.Interaction, fmt.Sprintf("Error occurred resetting player password: %v", err))
				if err != nil {
					log.Printf("Error sending message: %v", err)
				}
				return
			}

			channel, err := session.UserChannelCreate(member.User.ID)

			if err == nil {
				_, err = session.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Minecraft Server Night Pix",
							Description: "Your password has been reset.",
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
			if err != nil {
				log.Printf("Error sending message: %v", err)
			}

			_, err = session.FollowupMessageCreate(session.State.User.ID, interactionCreate.Interaction, true, &discordgo.WebhookParams{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Password reset",
						Description: "Successfully reset password",
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
						},
					},
				},
				Flags: 1 << 6,
			})
			if err != nil {
				log.Printf("Error sending message: %v", err)
				return
			}
		},
	},
}

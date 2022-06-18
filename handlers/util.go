package handlers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/minecraft-server-bot/connection"
	"log"
)

func createErrorEmbed(errorMessage string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Error",
		Description: errorMessage,
		Color:       ErrorEmbedColor,
	}
}
func interactionRespondError(session *discordgo.Session, interaction *discordgo.Interaction, errorMessage string) {
	err := session.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				createErrorEmbed(errorMessage),
			},
			Flags: 1 << 6,
		},
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
	}
}
func followupErrorMessageCreate(session *discordgo.Session, interaction *discordgo.Interaction, errorMessage string) {
	_, err := session.FollowupMessageCreate(session.State.User.ID, interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			createErrorEmbed(errorMessage),
		},
		Flags: 1 << 6,
	})
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func resetPasswordHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
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
		followupErrorMessageCreate(session, interactionCreate.Interaction, fmt.Sprintf("Error occurred getting player: %v", err))
		return
	}

	password, err := connection.ResetPlayerPassword(player.MinecraftNickname)
	if err != nil {
		log.Printf("Error resetting player password: %v", err)
		followupErrorMessageCreate(session, interactionCreate.Interaction, fmt.Sprintf("Error occurred resetting player password: %v", err))
		return
	}

	channel, messageErr := session.UserChannelCreate(member.User.ID)

	if messageErr == nil {
		_, messageErr = session.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
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
	if messageErr != nil {
		log.Printf("Error sending message: %v", err)
	}

	_, err = session.FollowupMessageCreate(session.State.User.ID, interactionCreate.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Password reset",
				Description: "Successfully reset password.",
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
	})
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return
	}
}

func updateWhitelistMessage(session *discordgo.Session) {
	channelSetting, err := connection.GetSetting("whitelistChannel")
	if err != nil {
		log.Printf("Error updating whitelist message: %v", err)
		return
	}
	messageSetting, err := connection.GetSetting("whitelistMessage")
	if err != nil {
		log.Printf("Error updating whitelist message: %v", err)
		return
	}

	players, err := connection.GetPlayers()
	if err != nil {
		log.Printf("Error updating whitelist message: %v", err)
		return
	}

	var fields []*discordgo.MessageEmbedField
	for _, player := range players {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  player.MinecraftNickname,
			Value: fmt.Sprintf("<@%v>", player.DiscordId),
		})
	}

	_, err = session.ChannelMessageEditComplex(&discordgo.MessageEdit{
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
		ID:      messageSetting.Value,
		Channel: channelSetting.Value,
	})
	if err != nil {
		log.Printf("Error updating whitelist message: %v", err)
		return
	}
}

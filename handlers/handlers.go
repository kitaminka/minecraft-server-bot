package handlers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/minecraft-server-bot/connection"
	"log"
)

var GuildId string

const (
	PrimaryEmbedColor = 9383347
	ErrorEmbedColor   = 13179932
)

var Handlers = []interface{}{
	func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
		switch interactionCreate.Type {
		case discordgo.InteractionApplicationCommand:
			Commands[interactionCreate.ApplicationCommandData().Name].Handler(session, interactionCreate)
		case discordgo.InteractionMessageComponent:
			Components[interactionCreate.MessageComponentData().CustomID].Handler(session, interactionCreate)
		case discordgo.InteractionModalSubmit:
			Modals[interactionCreate.ModalSubmitData().CustomID].Handler(session, interactionCreate)
		}
	},
	func(session *discordgo.Session, ready *discordgo.Ready) {
		updateWhitelistMessage(session)
	},
	func(session *discordgo.Session, guildMemberRemove *discordgo.GuildMemberRemove) {
		_, err := connection.UnregisterPlayer(guildMemberRemove.User.ID)
		if err != nil {
			log.Printf("Error unregistring player: %v", err)
			return
		}
		updateWhitelistMessage(session)
	},
}

func AddHandlers(session *discordgo.Session) {
	for _, value := range Handlers {
		session.AddHandler(value)
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

	user := interactionCreate.Member.User

	player, err := connection.GetPlayerByDiscord(user.ID)
	if err != nil {
		log.Printf("Error getting player: %v", err)
		followupErrorMessageCreate(session, interactionCreate.Interaction, fmt.Sprintf("Error occurred getting player: %v", err))
		return
	}

	password, err := connection.ResetMinecraftPlayerPassword(player.MinecraftNickname)
	if err != nil {
		log.Printf("Error resetting player password: %v", err)
		followupErrorMessageCreate(session, interactionCreate.Interaction, fmt.Sprintf("Error occurred resetting player password: %v", err))
		return
	}

	channel, messageErr := session.UserChannelCreate(user.ID)

	if messageErr == nil {
		_, messageErr = session.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Minecraft Server Night Pix",
					Description: "Your password has been reset.",
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

	_, err = session.FollowupMessageCreate(interactionCreate.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Password reset",
				Description: "Successfully reset password.",
				Color:       PrimaryEmbedColor,
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

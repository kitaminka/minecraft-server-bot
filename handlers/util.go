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
func interactionRespondError(session *discordgo.Session, interaction *discordgo.Interaction, errorMessage string) error {
	return session.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				createErrorEmbed(errorMessage),
			},
			Flags: 1 << 6,
		},
	})
}
func followupErrorMessageCreate(session *discordgo.Session, interaction *discordgo.Interaction, errorMessage string) (*discordgo.Message, error) {
	message, err := session.FollowupMessageCreate(session.State.User.ID, interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			createErrorEmbed(errorMessage),
		},
		Flags: 1 << 6,
	})
	if err != nil {
		return nil, err
	}
	return message, nil
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
					Components["reset-password"].MessageComponent,
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

package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/minecraft-server-bot/config"
	"log"
)

func createErrorEmbed(errorMessage string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Error",
		Description: errorMessage,
		Color:       config.Config.EmbedColors.Error,
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
		// TODO Remove logging, return error
		log.Printf("Error responding to interaction: %v", err)
	}
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

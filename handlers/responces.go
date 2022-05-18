package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/minecraft-server-bot/config"
	"log"
)

func interactionRespondError(session *discordgo.Session, interaction *discordgo.Interaction, errorMessage string) {
	err := session.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Error",
					Description: errorMessage,
					Color:       config.Config.EmbedColors.Error,
				},
			},
			Flags: 1 << 6,
		},
	})

	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
	}
}

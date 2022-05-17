package handlers

import (
	"github.com/bwmarrin/discordgo"
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
					Color:       9383347,
				},
			},
			Flags: 1 << 6,
		},
	})

	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
	}
}

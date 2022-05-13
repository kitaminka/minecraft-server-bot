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
					Title:       "Ошибка",
					Description: errorMessage,
					Color:       8523465,
				},
			},
			Flags: 1 << 6,
		},
	})

	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
	}
}
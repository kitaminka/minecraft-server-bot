package handlers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/server-bot/database"
	"log"
)

func GuildMemberAdd(discordSession *discordgo.Session, guildMemberAdd *discordgo.GuildMemberAdd) {
	guild, _ := discordSession.Guild(guildMemberAdd.GuildID)

	_, err := discordSession.ChannelMessageSendComplex("957693593693339698", &(discordgo.MessageSend{
		Content: fmt.Sprintf("<@%v>", guildMemberAdd.User.ID),
		Embeds: []*discordgo.MessageEmbed{{
			URL:         "",
			Title:       fmt.Sprintf("Welcome, %v!", guildMemberAdd.User.Username),
			Description: fmt.Sprintf("Welcome to the %v server!", guild.Name),
			Timestamp:   "",
			Color:       8523465,
			Footer:      nil,
			Image:       nil,
			Thumbnail:   nil,
			Video:       nil,
			Provider:    nil,
			Author:      nil,
			Fields:      nil,
		}},
		TTS: false,
		Components: []discordgo.MessageComponent{discordgo.ActionsRow{Components: []discordgo.MessageComponent{discordgo.Button{
			Label:    "Create profile",
			Style:    discordgo.PrimaryButton,
			Disabled: false,
			Emoji: discordgo.ComponentEmoji{
				Name: "➕",
			},
			URL:      "",
			CustomID: "create_profile",
		}}}},
		Files:           nil,
		AllowedMentions: nil,
		Reference:       nil,
		File:            nil,
		Embed:           nil,
	}))

	if err != nil {
		log.Printf("Error sending message: %v", err)
	}

	database.CreateNewMember(*guildMemberAdd.Member)
}

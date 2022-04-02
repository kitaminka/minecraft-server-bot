package handlers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/server-bot/config"
	"github.com/kitaminka/server-bot/db"
	"log"
)

func GuildMemberAdd(session *discordgo.Session, guildMemberAdd *discordgo.GuildMemberAdd) {
	guild, _ := session.Guild(guildMemberAdd.GuildID)

	_, err := session.ChannelMessageSendComplex(config.Config.Channels.WelcomeMessageChannel, &(discordgo.MessageSend{
		Content: fmt.Sprintf("<@%v>", guildMemberAdd.User.ID),
		Embeds: []*discordgo.MessageEmbed{
			{
				URL:         "",
				Title:       fmt.Sprintf("Добро пожаловать, **%v**!", guildMemberAdd.User.Username),
				Description: fmt.Sprintf("Вы зашли на сервер **%v**!", guild.Name),
				Timestamp:   "",
				Color:       8523465,
				Footer:      nil,
				Image:       nil,
				Thumbnail:   nil,
				Video:       nil,
				Provider:    nil,
				Author:      nil,
				Fields:      nil,
			},
		},
		TTS:             false,
		Components:      nil,
		Files:           nil,
		AllowedMentions: nil,
		Reference:       nil,
		File:            nil,
		Embed:           nil,
	}))

	if err != nil {
		log.Printf("Error sending message: %v", err)
	}

	db.CreateNewMember(*guildMemberAdd.Member)
}

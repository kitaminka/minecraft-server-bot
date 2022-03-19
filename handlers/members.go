package handlers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/server-bot/messages"
	"log"
)

func GuildMemberAdd(discordSession *discordgo.Session, guildMemberAdd *discordgo.GuildMemberAdd) {
	messages.DiscordMessages.WelcomeMessageEmbed.Description = fmt.Sprintf(messages.DiscordMessages.WelcomeMessageEmbed.Description, guildMemberAdd.User.Username)
	_, err := discordSession.ChannelMessageSendComplex("909856906582040629", &(discordgo.MessageSend{
		Content:         "Welcome!",
		Embeds:          []*discordgo.MessageEmbed{&messages.DiscordMessages.WelcomeMessageEmbed},
		TTS:             false,
		Components:      nil,
		Files:           nil,
		AllowedMentions: nil,
		Reference:       nil,
		File:            nil,
		Embed:           nil,
	}))
	if err != nil {
		log.Fatalf("Error sending welcome message: %v", err)
	}
}

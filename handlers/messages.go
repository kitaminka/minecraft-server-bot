package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/minecraft-server-bot/connection"
	"log"
	"strconv"
	"strings"
)

// Error messages
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
	_, err := session.FollowupMessageCreate(interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			createErrorEmbed(errorMessage),
		},
		Flags: 1 << 6,
	})
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

// Whitelist message
func createWhitelistEmbed() (*discordgo.MessageEmbed, error) {
	whitelistPlayers, err := connection.GetPlayerWhitelist()
	if err != nil {
		return &discordgo.MessageEmbed{}, err
	}
	whitelistPlayerCount, err := connection.GetWhitelistPlayerCount()
	if err != nil {
		return &discordgo.MessageEmbed{}, err
	}
	playerCount, err := connection.GetPlayerCount()
	if err != nil {
		return &discordgo.MessageEmbed{}, err
	}

	var whitelistPlayerString string
	if len(strings.Join(whitelistPlayers, ", ")) <= 1021 {
		whitelistPlayerString = strings.Join(whitelistPlayers, ", ")
	} else {
		whitelistPlayerString = strings.Join(whitelistPlayers, ", ")[:1021] + "..."
	}

	embed := &discordgo.MessageEmbed{
		Title: "Whitelist info",
		Color: PrimaryEmbedColor,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Total whitelist player count",
				Value: strconv.Itoa(whitelistPlayerCount),
			},
			{
				Name:  "Registered player count",
				Value: strconv.Itoa(playerCount),
			},
			{
				Name:  "Whitelist players",
				Value: whitelistPlayerString,
			},
		},
	}

	return embed, err
}

func updateWhitelistMessage(session *discordgo.Session) {
	channelSetting, err := connection.GetSetting(connection.WhitelistChannelSetting)
	if err != nil {
		log.Printf("Error updating whitelist message: %v", err)
		return
	}
	messageSetting, err := connection.GetSetting(connection.WhitelistMessageSetting)
	if err != nil {
		log.Printf("Error updating whitelist message: %v", err)
		return
	}

	embed, err := createWhitelistEmbed()
	if err != nil {
		log.Printf("Error creating whitelist message: %v", err)
		return
	}

	_, err = session.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Embeds: []*discordgo.MessageEmbed{
			embed,
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

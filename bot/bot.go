package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/server-bot/handlers"
	"log"
	"os"
	"os/signal"
)

func StartBot(token string) {
	discordSession, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	discordSession.AddHandler(handlers.GuildMemberAdd)

	discordSession.Identify.Intents = 1535

	err = discordSession.Open()
	if err != nil {
		log.Fatalf("Error opening Discord session: %v", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	discordSession.Close()
}

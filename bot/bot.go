package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/server-bot/handlers"
	"github.com/kitaminka/server-bot/util"
	"log"
	"os"
	"os/signal"
)

func StartBot(token string) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	session.AddHandler(handlers.GuildMemberAdd)
	session.AddHandler(handlers.InteractionCreate)

	session.Identify.Intents = util.Config.Intents

	err = session.Open()
	if err != nil {
		log.Fatalf("Error opening Discord session: %v", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	session.Close()
}

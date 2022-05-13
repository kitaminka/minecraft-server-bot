package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/server-bot/config"
	"github.com/kitaminka/server-bot/handlers"
	"log"
	"os"
	"os/signal"
)

func StartBot(token string) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Panicf("Error creating Discord session: %v", err)
	}

	handlers.AddHandlers(session)

	session.Identify.Intents = config.Config.Intents

	err = session.Open()
	if err != nil {
		log.Panicf("Error opening Discord session: %v", err)
	}

	handlers.CreateApplicationCommands(session)

	defer session.Close()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	// TODO Check if it works correctly

	if config.Config.RemoveApplicationCommands {
		handlers.RemoveApplicationCommands(session)
	}
}

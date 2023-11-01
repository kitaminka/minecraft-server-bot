package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/minecraft-server-bot/handlers"
	"log"
	"os"
	"os/signal"
)

const (
	Intents = 1535
)

func StartBot(token, guildId string) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Panicf("Error creating Discord session: %v", err)
	}

	handlers.AddHandlers(session)
	session.Identify.Intents = Intents

	err = session.Open()
	if err != nil {
		log.Panicf("Error opening Discord session: %v", err)
	}

	handlers.CreateApplicationCommands(session, guildId)
	log.Println("Bot is now running. Press CTRL-C to exit.")

	defer session.Close()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}

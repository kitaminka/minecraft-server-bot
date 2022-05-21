package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/minecraft-server-bot/config"
	"github.com/kitaminka/minecraft-server-bot/connection"
	"github.com/kitaminka/minecraft-server-bot/handlers"
	"log"
	"os"
	"os/signal"
)

func StartBot() {
	config.LoadConfig()
	config.LoadEnv()
	connection.ConnectMongo()
	connection.ConnectRcon(config.Config.Rcon.Address, config.Config.Rcon.Password)

	session, err := discordgo.New("Bot " + config.Config.Token)
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

	if config.Config.RemoveCommands {
		handlers.RemoveApplicationCommands(session)
	}
}

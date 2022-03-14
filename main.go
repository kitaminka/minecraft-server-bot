package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/kitaminka/server-bot/handlers"
	"log"
	"os"
	"os/signal"
)

var (
	Token string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	Token = os.Getenv("TOKEN")
}

func main() {
	discordSession, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	discordSession.AddHandler(handlers.Ready)

	discordSession.Identify.Intents = discordgo.IntentsAll

	err = discordSession.Open()
	if err != nil {
		log.Fatalf("Error opening Discord session: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	err = discordSession.Close()
	if err != nil {
		log.Fatalf("Error closing Discord session: %v", err)
	}
}

package main

import (
	"github.com/joho/godotenv"
	"github.com/kitaminka/server-bot/bot"
	"log"
	"os"
)

var Token string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	Token = os.Getenv("DISCORD_TOKEN")
}

func main() {
	bot.StartBot(Token)
}

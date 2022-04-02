package main

import (
	"github.com/joho/godotenv"
	"github.com/kitaminka/server-bot/bot"
	"github.com/kitaminka/server-bot/config"
	"github.com/kitaminka/server-bot/db"
	"log"
	"os"
)

var (
	Token    string
	MongoUri string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Panicf("Error loading .env file: %v", err)
	}
	Token = os.Getenv("DISCORD_TOKEN")
	MongoUri = os.Getenv("MONGODB_URI")
}

func main() {
	config.LoadConfig()
	db.Connect(MongoUri)
	bot.StartBot(Token)
}

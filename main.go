package main

import (
	"github.com/joho/godotenv"
	"github.com/kitaminka/server-bot/bot"
	"github.com/kitaminka/server-bot/db"
	"github.com/kitaminka/server-bot/util"
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
	util.LoadConfig()
	db.Connect(MongoUri)
	bot.StartBot(Token)
}

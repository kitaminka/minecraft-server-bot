package main

import (
	"github.com/joho/godotenv"
	"github.com/kitaminka/minecraft-server-bot/bot"
	"github.com/kitaminka/minecraft-server-bot/connection"
	"log"
	"os"
)

var (
	Token        string
	MongoUri     string
	RconAddress  string
	RconPassword string
	GuildId      string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Panicf("Error loading .env file: %v", err)
	}

	Token = os.Getenv("DISCORD_TOKEN")
	MongoUri = os.Getenv("MONGODB_URI")
	RconAddress = os.Getenv("RCON_ADDRESS")
	RconPassword = os.Getenv("RCON_PASSWORD")
	GuildId = os.Getenv("GUILD_ID")

	log.Print("Successfully loaded .env file")
}
func main() {
	connection.ImportConfiguration()
	connection.ConnectMongo(MongoUri)
	connection.ConnectRcon(RconAddress, RconPassword)
	bot.StartBot(Token, GuildId)
}

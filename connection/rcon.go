package connection

import (
	"crypto/rand"
	"fmt"
	"github.com/willroberts/minecraft-client"
	"log"
	"math/big"
	"strings"
)

// TODO Change all bool to error

var RconClient *minecraft.Client

func ConnectRcon(rconAddress, rconPassword string) {
	rconClient, err := minecraft.NewClient(rconAddress)
	if err != nil {
		log.Panicf("Error connecting to RCON: %v", err)
	}

	err = rconClient.Authenticate(rconPassword)
	if err != nil {
		log.Panicf("Error RCON authenticating: %v", err)
	}
	log.Print("Successfully connected to RCON")

	RconClient = rconClient
}
func RegisterPlayer(minecraftNickname string) (string, bool) {
	password := generatePassword()
	UnregisterPlayer(minecraftNickname)
	_, err := RconClient.SendCommand(fmt.Sprintf("nlogin register %v %v", minecraftNickname, password))
	if err != nil {
		log.Printf("Error sending command: %v", err)
		return "", false
	}
	return password, true
}
func UnregisterPlayer(minecraftNickname string) bool {
	_, err := RconClient.SendCommand(fmt.Sprintf("nlogin unregister %v", minecraftNickname))
	if err != nil {
		log.Printf("Error sending command: %v", err)
		return false
	}
	return true
}
func AddPlayerWhitelist(minecraftNickname string) bool {
	message, err := RconClient.SendCommand(fmt.Sprintf("whitelist add %v", minecraftNickname))
	if err != nil {
		log.Printf("Error sending command: %v", err)
		return false
	} else if message.Body == "Player is already whitelisted" {
		log.Printf("Player already exists: %v", minecraftNickname)
		return false
	}
	return true
}
func RemovePlayerWhitelist(minecraftNickname string) bool {
	_, err := RconClient.SendCommand(fmt.Sprintf("whitelist remove %v", minecraftNickname))
	if err != nil {
		log.Printf("Error sending command: %v", err)
		return false
	}
	return true
}
func GetPlayerWhitelist() []string {
	message, err := RconClient.SendCommand("whitelist list")
	if err != nil {
		log.Printf("Error sending command: %v", err)
		return nil
	}
	players := strings.Split(message.Body[33:], ", ")
	return players
}
func generatePassword() string {
	var password string

	for i := 0; i < 8; i++ {
		res, _ := rand.Int(rand.Reader, big.NewInt(61))
		password += string(Chars[res.Int64()])
	}

	return password
}

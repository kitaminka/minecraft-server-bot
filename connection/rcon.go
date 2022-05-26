package connection

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/willroberts/minecraft-client"
	"log"
	"math/big"
	"strings"
)

const passwordChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

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
func RegisterPlayer(minecraftNickname string) (string, error) {
	password := generatePassword()
	err := UnregisterPlayer(minecraftNickname)
	if err != nil {
		return "", err
	}
	_, err = RconClient.SendCommand(fmt.Sprintf("nlogin register %v %v", minecraftNickname, password))
	if err != nil {
		log.Printf("Error sending command: %v", err)
		return "", err
	}
	return password, nil
}
func UnregisterPlayer(minecraftNickname string) error {
	_, err := RconClient.SendCommand(fmt.Sprintf("nlogin unregister %v", minecraftNickname))
	if err != nil {
		return err
	}
	return nil
}
func AddPlayerWhitelist(minecraftNickname string) error {
	message, err := RconClient.SendCommand(fmt.Sprintf("whitelist add %v", minecraftNickname))
	if err != nil {
		log.Printf("Error sending command: %v", err)
		return err
	} else if message.Body == "Player is already whitelisted" {
		err = errors.New("player already exists")
		return err
	}
	return nil
}
func RemovePlayerWhitelist(minecraftNickname string) error {
	_, err := RconClient.SendCommand(fmt.Sprintf("whitelist remove %v", minecraftNickname))
	if err != nil {
		return err
	}
	return nil
}
func GetPlayerWhitelist() ([]string, error) {
	message, err := RconClient.SendCommand("whitelist list")
	if err != nil {
		log.Printf("Error sending command: %v", err)
		return nil, err
	}
	players := strings.Split(message.Body[33:], ", ")
	return players, nil
}
func generatePassword() string {
	var password string

	for i := 0; i < 8; i++ {
		res, _ := rand.Int(rand.Reader, big.NewInt(61))
		password += string(passwordChars[res.Int64()])
	}

	return password
}

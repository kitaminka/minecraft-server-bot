package connection

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

var Config Configuration

type Configuration struct {
	Commands struct {
		Register        string `json:"register"`
		Unregister      string `json:"unregister"`
		ChangePassword  string `json:"changePassword"`
		AddWhitelist    string `json:"addWhitelist"`
		RemoveWhitelist string `json:"removeWhitelist"`
		GetWhitelist    string `json:"getWhitelist"`
		Playtime        string `json:"playtime"`
	} `json:"commands"`
}

func ImportConfiguration() {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		log.Panicf("Error opening config.json: %v", err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &Config)
	if err != nil {
		log.Panicf("Error unmarshaling configuration: %v", err)
	}

	fmt.Println(Config)
}

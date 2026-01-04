package config

import (
	"encoding/json"
	"os"
)

type TwitchConfig struct {
	ReadAllChat bool   `json:"read_all_chat"`
	RedeemName  string `json:"redeem_name"`
}

type Config struct {
	TwitchConfig TwitchConfig `json:"twitch_config"`
}

var configFileName = "tts_config.json"

func ReadConfig() (Config, error) {
	openFile, err := os.ReadFile(configFileName)
	if err != nil {
		if os.IsNotExist(err) {
			_, err := os.Create(configFileName)
			if err != nil {
				return Config{}, err
			}
			return Config{}, nil
		} else {
			return Config{}, err
		}
	}

	var payload Config
	err = json.Unmarshal(openFile, &payload)
	if err != nil {
		return Config{}, err
	}
	return payload, nil
}

func SaveConfig(readAllChat bool, redeemName string) {
	var twitchConfigs TwitchConfig
	twitchConfigs.ReadAllChat = readAllChat
	twitchConfigs.RedeemName = redeemName

	var configs Config
	configs.TwitchConfig = twitchConfigs

	file, err := os.Create(configFileName)
	if err != nil {

	}

	defer file.Close()
	jsonData, err := json.Marshal(configs)
	if err != nil {

	}
	_, err = file.Write(jsonData)
	if err != nil {

	}
}

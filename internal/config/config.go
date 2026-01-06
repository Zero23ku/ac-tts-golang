package config

import (
	"encoding/json"
	"os"
	"strings"
)

type TwitchConfig struct {
	ReadAllChat bool   `json:"read_all_chat"`
	RedeemName  string `json:"redeem_name"`
}

type WhitelistConfig struct {
	RawWhitelist string `json:"raw_whitelist"`
}

type Config struct {
	TwitchConfig    TwitchConfig    `json:"twitch_config"`
	WhitelistConfig WhitelistConfig `json:"whitelist_config"`
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

func SaveConfig(readAllChat bool, redeemName string, whitelist []string) {
	var twitchConfigs TwitchConfig
	var whitelistconfs WhitelistConfig
	twitchConfigs.ReadAllChat = readAllChat
	twitchConfigs.RedeemName = redeemName
	whitelistconfs.RawWhitelist = strings.Join(whitelist, ",")

	var configs Config
	configs.TwitchConfig = twitchConfigs
	configs.WhitelistConfig = whitelistconfs

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

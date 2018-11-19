package config

import (
	"encoding/json"
	"os"
	"fmt"
)

type Chat struct {
	Id int64
	Name string
}

type Configuration struct {
	BannedNews  []string
	Chats     []Chat
	BotToken    string
}

func GetConfiguration() (Configuration, error) {
	file, _ := os.Open("/Users/wangchaoqun/Codebase/go/src/telegrambot/config/config_dev.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
		return configuration, err
	}

	return configuration, nil

}

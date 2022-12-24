package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"

	"github.com/pocket-bot/config"
	"github.com/pocket-bot/pkg/telegram"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("initialize config; error: %s", err.Error())
	}

	botApi, err := tgbotapi.NewBotAPI(viper.GetString("tlg_token"))
	if err != nil {
		log.Fatalf("initialize bot api; error: %s", err.Error())
	}
	botApi.Debug = true

	bot := telegram.NewBot(botApi)

	if err = bot.Start(); err != nil {
		log.Fatalf("start bot; error: %s", err.Error())
	}
}

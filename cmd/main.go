package main

import (
	"log"

	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
	"github.com/zhashkevych/go-pocket-sdk"
	"go.uber.org/zap"

	"github.com/pocket-bot/config"
	"github.com/pocket-bot/pkg/repository"
	"github.com/pocket-bot/pkg/repository/boltdb"
	"github.com/pocket-bot/pkg/server"
	"github.com/pocket-bot/pkg/telegram"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("create zap logger; error: %s", err.Error())
	}

	if err := config.InitConfig(); err != nil {
		log.Fatalf("initialize config; error: %s", err.Error())
	}

	db, err := initBolt()
	if err != nil {
		log.Fatalf("initialize bolt db; error: %s", err.Error())
	}
	storage := boltdb.NewStorage(db)

	botApi, err := tgbotapi.NewBotAPI(viper.GetString("tlg_token"))
	if err != nil {
		log.Fatalf("initialize bot api; error: %s", err.Error())
	}
	botApi.Debug = true

	pocketClient, err := pocket.NewClient(viper.GetString("pocket_key"))
	if err != nil {
		log.Fatalf("initalize pocket client; error: %s", err.Error())
	}

	bot := telegram.NewBot(botApi, pocketClient, storage, "http://localhost:8000/", logger)
	go func() {
		if err = bot.Start(); err != nil {
			log.Fatalf("start bot; error: %s", err.Error())
		}
	}()

	httpServer := server.NewHttpServer(storage, pocketClient, logger, "https://web.telegram.org/k/#@livert21_bot")
	if err := httpServer.Start(":8000"); err != nil {
		log.Fatalf("start server; error: %s", err.Error())
	}
}

func initBolt() (*bolt.DB, error) {
	db, err := bolt.Open("bot.db", 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Batch(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessTokens))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(repository.RequestTokens))
		return err
	}); err != nil {
		return nil, err
	}

	return db, nil
}

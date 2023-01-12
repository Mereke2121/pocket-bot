package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"github.com/zhashkevych/go-pocket-sdk"
	"go.uber.org/zap"

	"github.com/pocket-bot/pkg/repository"
	"github.com/pocket-bot/pkg/repository/boltdb"
)

type Bot struct {
	bot          *tgbotapi.BotAPI
	pocketClient *pocket.Client
	storage      repository.TokenStorage
	redirectURL  string
	logger       *zap.Logger
}

func NewBot(bot *tgbotapi.BotAPI, pocketClient *pocket.Client, storage *boltdb.TokenStorage, redirectURL string, logger *zap.Logger) *Bot {
	return &Bot{
		bot:          bot,
		pocketClient: pocketClient,
		storage:      storage,
		redirectURL:  redirectURL,
		logger:       logger,
	}
}

func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.bot.GetUpdatesChan(u)
	if err != nil {
		return errors.Wrap(err, "get updates channel")
	}

	for update := range updates {
		if update.Message == nil { // If we didn't get a message
			continue
		}
		if update.Message.IsCommand() {
			b.handleCommand(*update.Message)
		} else {
			b.handleMessage(*update.Message)
		}
	}
	return nil
}

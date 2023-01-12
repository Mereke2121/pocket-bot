package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/pocket-bot/pkg/repository"
)

func (b *Bot) handleCommand(message tgbotapi.Message) {
	var text string
	switch message.Command() {
	case "start":
		_, err := b.storage.Get(int(message.Chat.ID), repository.AccessTokens)
		if err != nil {
			text = b.initAuthorization(message)
		} else {
			text = "Вы уже авторизованы"
		}
	default:
		text = "Неизвестная команда :("
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	b.bot.Send(msg)
}

func (b *Bot) handleMessage(message tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	b.bot.Send(msg)
}

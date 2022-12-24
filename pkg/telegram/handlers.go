package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
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
}

func (b *Bot) handleCommand(message tgbotapi.Message) {
	var text string
	switch message.Command() {
	case "start":
		text = "Ты ввел команду старт"
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

package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"

	"github.com/pocket-bot/pkg/repository"
)

func (b *Bot) initAuthorization(message tgbotapi.Message) string {
	redirectUrl := fmt.Sprintf("%s?chat_id=%v", b.redirectURL, message.Chat.ID)
	requestToken, _ := b.pocketClient.GetRequestToken(context.Background(), redirectUrl)
	err := b.storage.Save(int(message.Chat.ID), requestToken, repository.RequestTokens)
	if err != nil {
		b.logger.Error("save request token in storage", zap.Error(err))
	}
	authURL, _ := b.pocketClient.GetAuthorizationURL(requestToken, redirectUrl)
	text := fmt.Sprintf("Привет! Чтобы сохранять ссылки в своем Pocket аккаунте, для начала тебе необходимо дать мне на это доступ. Для этого переходи по ссылке: %s", authURL)
	return text
}

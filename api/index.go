package api

import (
	"github.com/41x3n/Xom/shared"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type rootHandler struct {
	env      *shared.Env
	telegram shared.TelegramService
}

func (h *rootHandler) HandleMessages(update tgbotapi.Update) error {
	tg := h.telegram.GetAPI()
	// echo back the message
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	_, err := tg.Send(msg)
	if err != nil {
		return errors.Wrap(err, "error sending message")
	}
	return nil
}

func NewRootHandler(env *shared.Env, telegram shared.TelegramService) shared.RootHandlerInterface {
	return &rootHandler{env: env, telegram: telegram}
}

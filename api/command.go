package api

import (
	"github.com/41x3n/Xom/api/controller"
	"github.com/41x3n/Xom/core/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *rootHandler) HandleStartCommand(user *domain.User, message *tgbotapi.Message) error {
	telegramAPI := h.telegram.GetAPI()

	sc := controller.NewStartController(telegramAPI)

	err := sc.HandleStartCommand(user, message)

	return err
}

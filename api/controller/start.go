package controller

import (
	"fmt"
	"github.com/41x3n/Xom/core/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StartController struct {
	TelegramAPI *tgbotapi.BotAPI
}

func NewStartController(telegramAPI *tgbotapi.BotAPI) *StartController {
	return &StartController{TelegramAPI: telegramAPI}
}

func (sc *StartController) HandleStartCommand(user *domain.User, message *tgbotapi.Message) error {
	var msgText string
	if user.UserName != nil {
		msgText = fmt.Sprintf("Howdy, %s! Welcome aboard the Xom express, your personal assistant for all your file conversion needs. Buckle up and enjoy the ride!", *user.UserName)
	} else {
		msgText = "Hello! Welcome to Xom, your butler to handle all your file conversion help."
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	_, err := sc.TelegramAPI.Send(msg)
	return err
}

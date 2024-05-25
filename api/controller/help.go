package controller

import (
	"github.com/41x3n/Xom/core/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HelpController struct {
	TelegramAPI *tgbotapi.BotAPI
}

func NewHelpController(telegramAPI *tgbotapi.BotAPI) *HelpController {
	return &HelpController{TelegramAPI: telegramAPI}
}

func (hc *HelpController) HandleHelpCommand(user *domain.User, message *tgbotapi.Message) error {
	msgText := "Xom is your personal assistant for all your file conversion needs. You can convert files from one format to another, and also extract text from images. Here are the commands you can use:\n\n"
	msgText += "/start - Start the bot\n"
	msgText += "/help - Get help\n"
	msgText += "/about - About Xom\n"
	msgText += "Or simply send a file to get started!"
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	_, err := hc.TelegramAPI.Send(msg)
	return err
}

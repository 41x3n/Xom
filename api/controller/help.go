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
	msgText := "Need a little guidance? Xom's got your back! 🌟\n\n"
	msgText += "🔄 **Supported Conversions:**\n\n"
	msgText += "🎵 **Audio Files:** Convert your tunes to any of these formats: mp4, mp3, wav, flac, ogg, aac, wma, m4a.\n\n"
	msgText += "📷 **Image Files:** Turn your images into: jpg, jpeg, png, gif, pdf, webp, bmp, tif, tiff, ico, avif.\n\n"
	msgText += "Just send me a file, and I’ll handle the rest. Let’s make those files fit your needs, effortlessly! 🎉"
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	_, err := hc.TelegramAPI.Send(msg)
	return err
}

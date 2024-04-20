package controller

import (
	"fmt"
	"github.com/41x3n/Xom/core/domain"
	"github.com/41x3n/Xom/shared"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackController struct {
	cu          domain.CallbackUseCase
	TelegramAPI *tgbotapi.BotAPI
}

func NewCallbackController(cu domain.CallbackUseCase,
	telegramAPI *tgbotapi.BotAPI) *CallbackController {
	return &CallbackController{
		cu:          cu,
		TelegramAPI: telegramAPI}
}

func (cc *CallbackController) HandleCallback(callback *tgbotapi.CallbackQuery) error {
	photoId, command, err := cc.cu.GetFileIDAndCommand(callback)
	if err != nil {
		return err
	}

	photoCommandString := fmt.Sprintf("%v", shared.PhotoCommand)
	if command == photoCommandString {
		photo, err := cc.cu.GetPhotoByID(photoId)
		if err != nil {
			return err
		}
		fmt.Println(photo)
	}

	return nil
}

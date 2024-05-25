package controller

import (
	"github.com/41x3n/Xom/core/domain"
	"github.com/41x3n/Xom/shared"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PhotoController struct {
	pu          domain.PhotoUseCase
	TelegramAPI *tgbotapi.BotAPI
}

func NewPhotoController(pu domain.PhotoUseCase, telegramAPI *tgbotapi.BotAPI) *PhotoController {
	return &PhotoController{
		pu:          pu,
		TelegramAPI: telegramAPI}
}

func (pc *PhotoController) HandlePhoto(user *domain.User, message *tgbotapi.Message) error {
	var photo *domain.Photo
	fileID, fileType, err := pc.pu.GetFileIDAndType(message)
	if err != nil {
		return err
	}

	var messageID = message.MessageID

	photo, err = pc.pu.SavePhotoId(user, fileID, fileType, messageID)
	if err != nil {
		return err
	}

	buttonRows := pc.pu.GenerateConvertOptions(fileType, photo.ID)
	if len(buttonRows) == 0 {
		return shared.ErrInvalidFile
	}

	keyboard := pc.pu.GenerateKeyboardMarkup(buttonRows)
	msg := pc.pu.GenerateMessage(fileType, message, keyboard)

	_, err = pc.TelegramAPI.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

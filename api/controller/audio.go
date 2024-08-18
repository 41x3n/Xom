package controller

import (
	"github.com/41x3n/Xom/core/domain"
	"github.com/41x3n/Xom/shared"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AudioController struct {
	usecase     domain.AudioUseCase
	TelegramAPI *tgbotapi.BotAPI
}

func NewAudioController(usecase domain.AudioUseCase, telegramAPI *tgbotapi.BotAPI) *AudioController {
	return &AudioController{
		usecase:     usecase,
		TelegramAPI: telegramAPI,
	}
}

func (ac *AudioController) HandleAudio(user *domain.User, message *tgbotapi.Message) error {
	fileID, fileType, errFileType := ac.usecase.GetFileIDAndType(message)
	if errFileType != nil {
		return errFileType
	}

	audio, errOnSave := ac.usecase.SaveAudioId(user, fileID, fileType, message.MessageID)
	if errOnSave != nil {
		return errOnSave
	}

	buttonRows := ac.usecase.GenerateConvertOptions(fileType, audio.ID)
	if len(buttonRows) == 0 {
		return shared.ErrInvalidFile
	}

	keyboard := ac.usecase.GenerateKeyboardMarkup(buttonRows)
	msg := ac.usecase.GenerateMessage(fileType, message, keyboard)

	_, err := ac.TelegramAPI.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

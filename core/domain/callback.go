package domain

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type CallbackUseCase interface {
	GetFileIDAndCommand(callback *tgbotapi.CallbackQuery) (string, string,
		string, error)
	GetPhotoByID(fileID string) (*Photo, error)
	MarkPhotoAsPreparing(photo *Photo) error
}

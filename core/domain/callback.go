package domain

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type CallbackUseCase interface {
	GetFileIDAndCommand(callback *tgbotapi.CallbackQuery) (string, string,
		string, error)
	GetPhotoByID(fileID string) (*Photo, error)
	GetAudioByID(fileID string) (*Audio, error)
	MarkPhotoAsPreparing(photo *Photo) error
	MarkAudioAsPreparing(audio *Audio) error
}

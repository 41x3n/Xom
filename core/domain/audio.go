package domain

import (
	"context"

	"github.com/41x3n/Xom/shared"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

var AudioFileTypeArray = []string{
	"mp4",
	"mp3",
	"wav",
	"flac",
	"ogg",
	"aac",
	"wma",
	"m4a",
}

type Audio struct {
	gorm.Model
	ID             int64         `gorm:"primary_key;autoIncrement"`
	UserTelegramID int64         `gorm:"not null;foreignKey:TelegramID"`
	FileID         string        `gorm:"not null"`
	FileType       string        `gorm:"not null"`
	Status         shared.Status `gorm:"not null;default:initiated;check:status IN ('initiated', 'preparing', 'processing', 'completed', 'failed')"`
	ConvertTo      string        `gorm:"not null;default:mp3"`
	MessageID      int64         `gorm:"default:null"`
}

type AudioRepository interface {
	Create(c context.Context, audio *Audio) error
	FetchByUser(c context.Context, user *User) ([]Audio, error)
	GetByFileID(c context.Context, fileID string) (Audio, error)
	GetByID(c context.Context, id int64) (*Audio, error)
	UpdateStatus(c context.Context, audio *Audio) error
}

type AudioUseCase interface {
	GetFileIDAndType(*tgbotapi.Message) (string, string, error)
	SaveAudioId(user *User, fileID, fileType string, messageID int) (*Audio, error)
	GenerateConvertOptions(fileType string, audioID int64) [][]tgbotapi.InlineKeyboardButton
	GenerateKeyboardMarkup(buttonRows [][]tgbotapi.InlineKeyboardButton) tgbotapi.InlineKeyboardMarkup
	GenerateMessage(fileType string, message *tgbotapi.Message, keyboard tgbotapi.InlineKeyboardMarkup) tgbotapi.MessageConfig
	ValidateIfAudioReadyToBeConverted(ID int64) (*Audio, error)
	UpdateAudioStatus(audio *Audio, status shared.Status) error
}

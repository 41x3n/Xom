package domain

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type Status string

const (
	Initiated  Status = "initiated"
	Processing Status = "processing"
	Completed  Status = "completed"
	Failed     Status = "failed"
)

var ErrInvalidFile = errors.New("invalid file")
var ErrInvalidCallbackData = errors.New("invalid callback data")

var FileTypeArray = []string{
	"jpg",
	"jpeg",
	"png",
	"gif",
	"pdf",
	"webp",
}

type Photo struct {
	gorm.Model
	ID             int64  `gorm:"primary_key;autoIncrement"`
	UserTelegramID int64  `gorm:"not null;foreignKey:TelegramID"`
	FileID         string `gorm:"not null"`
	FileType       string `gorm:"not null"`
	Status         Status `gorm:"not null;default:initiated;check:status IN ('initiated', 'processing', 'completed', 'failed')"`
	ConvertTo      string `gorm:"not null;default:jpg"`
}

type PhotoRepository interface {
	Create(c context.Context, photo *Photo) error
	FetchByUser(c context.Context, user *User) ([]Photo, error)
	GetByFileID(c context.Context, fileID string) (Photo, error)
	GetByID(c context.Context, id int64) (*Photo, error)
}

type PhotoUseCase interface {
	GetFileIDAndType(*tgbotapi.Message) (string, string, error)
	SavePhotoId(user *User, fileID, fileType string) (*Photo, error)
	GenerateConvertOptions(fileType string, photoID int64) [][]tgbotapi.InlineKeyboardButton
	GenerateKeyboardMarkup(buttonRows [][]tgbotapi.InlineKeyboardButton) tgbotapi.InlineKeyboardMarkup
	GenerateMessage(fileType string, message *tgbotapi.Message, keyboard tgbotapi.InlineKeyboardMarkup) tgbotapi.MessageConfig
}

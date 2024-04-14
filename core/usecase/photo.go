package usecase

import (
	"context"
	"fmt"
	"github.com/41x3n/Xom/core/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"time"
)

type photoUsecase struct {
	photoRepo      domain.PhotoRepository
	contextTimeout time.Duration
}

func NewPhotoUsecase(pr domain.PhotoRepository, timeout time.Duration) domain.PhotoUseCase {
	return &photoUsecase{
		photoRepo:      pr,
		contextTimeout: timeout,
	}
}

func (pu *photoUsecase) GetFileIDAndType(message *tgbotapi.Message) (string, string, error) {
	var fileID string
	var fileType string

	if message.Photo != nil && len(message.Photo) > 0 {
		fileID = message.Photo[3].FileID
		fileType = "jpg"
	}

	if message.Document != nil {
		fileID = message.Document.FileID
		mimeTypeParts := strings.Split(message.Document.MimeType, "/")
		if len(mimeTypeParts) > 1 {
			fileType = mimeTypeParts[1]
		} else {
			fileType = ""
		}
	}

	if fileID == "" || fileType == "" {
		return "", "", domain.ErrInvalidFile
	}

	return fileID, fileType, nil
}

func (pu *photoUsecase) SavePhotoId(user *domain.User, fileID,
	fileType string) (*domain.Photo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), pu.contextTimeout)
	defer cancel()

	photoStruct := &domain.Photo{
		UserTelegramID: user.TelegramID,
		FileID:         fileID,
		FileType:       fileType,
		Status:         domain.Initiated,
	}

	err := pu.photoRepo.Create(ctx, photoStruct)

	if err != nil {
		return nil, err
	}

	return photoStruct, nil
}

func (pu *photoUsecase) GenerateConvertOptions(
	ignoredFileType string) [][]tgbotapi.InlineKeyboardButton {
	ignoredBtn := strings.ToLower(ignoredFileType)

	var buttons [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	for i, fileType := range domain.FileTypeArray {
		if fileType != ignoredBtn {
			button := tgbotapi.NewInlineKeyboardButtonData(strings.ToUpper(fileType), fileType)
			row = append(row, button)
		}
		if len(row) == 2 || i == len(domain.FileTypeArray)-1 {
			buttons = append(buttons, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}

	return buttons
}

func (pu *photoUsecase) GenerateKeyboardMarkup(
	buttonRows [][]tgbotapi.InlineKeyboardButton) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, buttonRow := range buttonRows {
		rows = append(rows, buttonRow)
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	return keyboard
}

func (pu *photoUsecase) GenerateMessage(fileType string, message *tgbotapi.Message,
	keyboard tgbotapi.
		InlineKeyboardMarkup) tgbotapi.MessageConfig {
	msgText := fmt.Sprintf("You have uploaded a %s file. "+
		"Please choose which format you want to convert.",
		strings.ToUpper(fileType))
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	msg.ReplyMarkup = keyboard

	return msg
}

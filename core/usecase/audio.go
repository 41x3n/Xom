package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/41x3n/Xom/shared"

	"github.com/41x3n/Xom/core/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type audioUsecase struct {
	audioRepo      domain.AudioRepository
	contextTimeout time.Duration
}

func NewAudioUsecase(ar domain.AudioRepository, timeout time.Duration) domain.AudioUseCase {
	return &audioUsecase{
		audioRepo:      ar,
		contextTimeout: timeout,
	}
}

func (u *audioUsecase) GetFileIDAndType(message *tgbotapi.Message) (string, string, error) {
	var fileID string
	var fileType string
	var mimeTypeParts []string

	if message.Audio != nil {
		fileID = message.Audio.FileID
		mimeTypeParts = strings.Split(message.Audio.MimeType, "/")
	}
	if message.Voice != nil {
		fileID = message.Voice.FileID
		mimeTypeParts = strings.Split(message.Voice.MimeType, "/")
	}

	if len(mimeTypeParts) > 1 {
		fileType = mimeTypeParts[1]
	} else {
		fileType = ""
	}

	if fileID == "" || fileType == "" {
		return "", "", shared.ErrInvalidFile
	}

	if fileType == "mpeg" {
		fileType = "mp3"
	}

	return fileID, fileType, nil
}

func (u *audioUsecase) SaveAudioId(user *domain.User, fileID, fileType string, messageID int) (*domain.Audio, error) {
	audio := &domain.Audio{
		UserTelegramID: user.TelegramID,
		FileID:         fileID,
		FileType:       fileType,
		MessageID:      int64(messageID),
	}
	ctx, cancel := context.WithTimeout(context.Background(), u.contextTimeout)
	defer cancel()
	err := u.audioRepo.Create(ctx, audio)
	return audio, err
}

func (u *audioUsecase) GenerateConvertOptions(ignoredFileType string, audioID int64) [][]tgbotapi.InlineKeyboardButton {
	ignoredBtn := strings.ToLower(ignoredFileType)

	var buttonRows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	for i, fileType := range domain.AudioFileTypeArray {
		if fileType == ignoredBtn {
			continue
		}
		btnLabel := strings.ToUpper(fileType)
		data := fmt.Sprintf("%d-%s-%s", audioID, shared.AudioCommand, fileType)
		btn := tgbotapi.NewInlineKeyboardButtonData(btnLabel, data)
		row = append(row, btn)

		if len(row) == 2 || i == len(domain.AudioFileTypeArray)-1 {
			buttonRows = append(buttonRows, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}

	return buttonRows
}

func (u *audioUsecase) GenerateKeyboardMarkup(buttonRows [][]tgbotapi.InlineKeyboardButton) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, buttonRow := range buttonRows {
		rows = append(rows, buttonRow)
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	return keyboard
}

func (u *audioUsecase) GenerateMessage(fileType string, message *tgbotapi.Message, keyboard tgbotapi.InlineKeyboardMarkup) tgbotapi.MessageConfig {
	msgText := fmt.Sprintf("%s file detected! "+
		"🎉 Now, choose your desired format for conversion.",
		strings.ToUpper(fileType))
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	msg.ReplyMarkup = keyboard

	return msg
}

func (u *audioUsecase) ValidateIfAudioReadyToBeConverted(ID int64) (*domain.Audio, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.contextTimeout)
	defer cancel()
	return u.audioRepo.GetByID(ctx, ID)
}

func (u *audioUsecase) UpdateAudioStatus(audio *domain.Audio, status shared.Status) error {
	ctx, cancel := context.WithTimeout(context.Background(), u.contextTimeout)
	defer cancel()
	audio.Status = status
	return u.audioRepo.UpdateStatus(ctx, audio)
}

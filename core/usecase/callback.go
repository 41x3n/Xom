package usecase

import (
	"context"
	"github.com/41x3n/Xom/core/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
	"time"
)

type CallbackUseCase struct {
	photoRepo      domain.PhotoRepository
	contextTimeout time.Duration
}

func NewCallbackUseCase(pr domain.PhotoRepository, timeout time.Duration) domain.CallbackUseCase {
	return &CallbackUseCase{
		photoRepo:      pr,
		contextTimeout: timeout,
	}
}

func (cu *CallbackUseCase) GetFileIDAndCommand(callback *tgbotapi.CallbackQuery) (string, string, error) {
	data := callback.Data
	if data == "" {
		return "", "", domain.ErrInvalidCallbackData
	}

	dataParts := strings.Split(data, "-")
	if len(dataParts) < 2 {
		return "", "", domain.ErrInvalidCallbackData

	}

	return dataParts[0], dataParts[1], nil
}

func (cu *CallbackUseCase) GetPhotoByID(photoId string) (*domain.Photo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cu.contextTimeout)
	defer cancel()

	id, parseErr := strconv.ParseInt(photoId, 10, 8)

	if parseErr != nil {
		return nil, parseErr
	}

	photo, err := cu.photoRepo.GetByID(ctx, id)

	return photo, err
}

func (cu *CallbackUseCase) MarkPhotoAsProcessing(photo *domain.Photo) error {
	ctx, cancel := context.WithTimeout(context.Background(), cu.contextTimeout)
	defer cancel()

	photo.Status = domain.Processing
	err := cu.photoRepo.UpdateStatus(ctx, photo)

	if err != nil {
		return err
	}

	return nil
}

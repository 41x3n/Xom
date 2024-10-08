package usecase

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/41x3n/Xom/core/domain"
	"github.com/41x3n/Xom/shared"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackUseCase struct {
	photoRepo      domain.PhotoRepository
	audioRepo      domain.AudioRepository
	contextTimeout time.Duration
}

func NewCallbackUseCase(pr domain.PhotoRepository, ar domain.AudioRepository, timeout time.Duration) domain.CallbackUseCase {
	return &CallbackUseCase{
		photoRepo:      pr,
		audioRepo:      ar,
		contextTimeout: timeout,
	}
}

func (cu *CallbackUseCase) GetFileIDAndCommand(callback *tgbotapi.
	CallbackQuery) (string, string, string, error) {
	data := callback.Data
	if data == "" {
		return "", "", "", shared.ErrInvalidCallbackData
	}

	dataParts := strings.Split(data, "-")
	if len(dataParts) < 3 {
		return "", "", "", shared.ErrInvalidCallbackData

	}

	return dataParts[0], dataParts[1], dataParts[2], nil
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

func (cu *CallbackUseCase) GetAudioByID(audioId string) (*domain.Audio, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cu.contextTimeout)
	defer cancel()

	id, parseErr := strconv.ParseInt(audioId, 10, 8)

	if parseErr != nil {
		return nil, parseErr
	}

	audio, err := cu.audioRepo.GetByID(ctx, id)

	return audio, err
}

func (cu *CallbackUseCase) MarkPhotoAsPreparing(photo *domain.Photo) error {
	ctx, cancel := context.WithTimeout(context.Background(), cu.contextTimeout)
	defer cancel()

	photo.Status = shared.Preparing
	err := cu.photoRepo.UpdateStatus(ctx, photo)

	if err != nil {
		return err
	}

	return nil
}

func (cu *CallbackUseCase) MarkAudioAsPreparing(audio *domain.Audio) error {
	ctx, cancel := context.WithTimeout(context.Background(), cu.contextTimeout)
	defer cancel()

	audio.Status = shared.Preparing
	err := cu.audioRepo.UpdateStatus(ctx, audio)

	if err != nil {
		return err
	}

	return nil
}

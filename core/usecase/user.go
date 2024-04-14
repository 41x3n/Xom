package usecase

import (
	"context"
	"fmt"
	"github.com/41x3n/Xom/core/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

type userUsecase struct {
	userRepo       domain.UserRepository
	contextTimeout time.Duration
}

func NewUserUsecase(ur domain.UserRepository, timeout time.Duration) domain.UserUseCase {
	return &userUsecase{
		userRepo:       ur,
		contextTimeout: timeout,
	}
}

func (uu *userUsecase) VerifyUser(tgUser tgbotapi.User) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), uu.contextTimeout)
	defer cancel()
	fmt.Println("Verifying user", tgUser.ID)

	userStruct := &domain.User{
		TelegramID:   tgUser.ID,
		IsBot:        tgUser.IsBot,
		FirstName:    tgUser.FirstName,
		LastName:     &tgUser.LastName,
		UserName:     &tgUser.UserName,
		LanguageCode: &tgUser.LanguageCode,
		IsActive:     true,
	}

	user, err := uu.userRepo.GetOrCreateByUserTelegramID(ctx, userStruct)
	return user, err
}

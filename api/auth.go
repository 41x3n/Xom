package api

import (
	"github.com/41x3n/Xom/core/domain"
	"github.com/41x3n/Xom/core/repository"
	"github.com/41x3n/Xom/core/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"time"
)

func checkIfUserExists(tgUser *tgbotapi.User, db *gorm.DB) (*domain.User,
	error) {
	ur := repository.NewUserRepository(db, domain.TableUser)
	uc := usecase.NewUserUsecase(ur, 50*time.Second)

	// Check if user exists
	user, err := uc.VerifyUser(*tgUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}

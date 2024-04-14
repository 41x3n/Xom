package domain

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"gorm.io/gorm"
)

const (
	TableUser = "users"
)

type User struct {
	gorm.Model
	TelegramID   int64  `gorm:"primaryKey;unique"`
	IsBot        bool   `gorm:"not null"`
	FirstName    string `gorm:"not null"`
	LastName     *string
	UserName     *string
	LanguageCode *string
	IsActive     bool `gorm:"default:true"`
	//Photos       []Photo `gorm:"foreignKey:UserTelegramID;references:TelegramID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type UserRepository interface {
	Create(c context.Context, user *User) error
	Fetch(c context.Context) ([]User, error)
	GetByUserName(c context.Context, username string) (User, error)
	GetByID(c context.Context, telegramID int64) (User, error)
	GetOrCreateByUserTelegramID(c context.Context, user *User) (*User, error)
}

type UserUseCase interface {
	VerifyUser(tgUser tgbotapi.User) (*User, error)
}

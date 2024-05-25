package repository

import (
	"context"

	"github.com/41x3n/Xom/core/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	database *gorm.DB
	table    string
}

func NewUserRepository(db *gorm.DB, table string) domain.UserRepository {
	return &userRepository{
		database: db,
		table:    table,
	}
}

func (ur *userRepository) Create(ctx context.Context, user *domain.User) error {
	result := ur.database.WithContext(ctx).Create(user)
	return result.Error
}

func (ur *userRepository) Fetch(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	result := ur.database.WithContext(ctx).Find(&users)
	return users, result.Error
}

func (ur *userRepository) GetByUserName(ctx context.Context, username string) (domain.User, error) {
	var user domain.User
	result := ur.database.WithContext(ctx).Where("user_name = ?", username).First(&user)
	return user, result.Error
}

func (ur *userRepository) GetByID(ctx context.Context, telegramID int64) (domain.User, error) {
	var user domain.User
	result := ur.database.WithContext(ctx).Where("telegram_id = ?", telegramID).First(&user)
	return user, result.Error
}

func (ur *userRepository) GetOrCreateByUserTelegramID(ctx context.Context,
	user *domain.User) (*domain.User, error) {
	err := ur.database.WithContext(ctx).Table(ur.table).Where("telegram_id = ?", user.TelegramID).FirstOrCreate(user).Error
	return user, err
}

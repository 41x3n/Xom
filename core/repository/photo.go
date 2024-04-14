package repository

import (
	"context"
	"github.com/41x3n/Xom/core/domain"
	"gorm.io/gorm"
)

type photoRepository struct {
	database *gorm.DB
	table    string
}

func NewPhotoRepository(db *gorm.DB, table string) domain.PhotoRepository {
	return &photoRepository{
		database: db,
		table:    table,
	}
}

func (pr *photoRepository) Create(ctx context.Context, photo *domain.Photo) error {
	result := pr.database.Create(photo)
	return result.Error
}

func (pr *photoRepository) FetchByUser(ctx context.Context, user *domain.User) ([]domain.Photo, error) {
	var photos []domain.Photo
	result := pr.database.Where("user_telegram_id = ?", user.TelegramID).Find(&photos)
	return photos, result.Error
}

func (pr *photoRepository) GetByFileID(ctx context.Context, fileID string) (domain.Photo, error) {
	var photo domain.Photo
	result := pr.database.Where("file_id = ?", fileID).First(&photo)
	return photo, result.Error
}

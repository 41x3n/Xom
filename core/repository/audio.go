package repository

import (
	"context"

	"github.com/41x3n/Xom/core/domain"
	"gorm.io/gorm"
)

type audioRepository struct {
	database *gorm.DB
	table    string
}

func NewAudioRepository(db *gorm.DB, table string) domain.AudioRepository {
	return &audioRepository{
		database: db,
		table:    table,
	}
}

func (r *audioRepository) Create(c context.Context, audio *domain.Audio) error {
	result := r.database.WithContext(c).Create(audio)
	return result.Error
}

func (r *audioRepository) FetchByUser(c context.Context, user *domain.User) ([]domain.Audio, error) {
	var audios []domain.Audio
	err := r.database.WithContext(c).Where("user_telegram_id = ?", user.TelegramID).Find(&audios).Error
	return audios, err
}

func (r *audioRepository) GetByFileID(c context.Context, fileID string) (domain.Audio, error) {
	var audio domain.Audio
	err := r.database.WithContext(c).Where("file_id = ?", fileID).First(&audio).Error
	return audio, err
}

func (r *audioRepository) GetByID(c context.Context, id int64) (*domain.Audio, error) {
	var audio domain.Audio
	err := r.database.WithContext(c).First(&audio, id).Error
	return &audio, err
}

func (r *audioRepository) UpdateStatus(c context.Context, audio *domain.Audio) error {
	return r.database.WithContext(c).Model(audio).Updates(audio).Error
}

package service

import (
	"github.com/41x3n/Xom/shared"
	"gorm.io/gorm"
	"log"
)

type converter struct {
	env      *shared.Env
	telegram shared.TelegramService
	db       *gorm.DB
}

func (c *converter) HandleFiles(payload *shared.RabbitMQPayload) error {
	var err error

	switch payload.Command {
	case shared.PhotoCommand:
		err = HandlePhotos(payload.ID)
	default:
		log.Printf("Unknown command: %s", payload.Command)
	}

	return err
}

func NewFFMPEGService(env *shared.Env, telegram shared.TelegramService, db *gorm.DB) shared.FFMPEGService {
	return &converter{
		env:      env,
		telegram: telegram,
		db:       db,
	}
}

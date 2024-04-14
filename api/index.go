package api

import (
	"github.com/41x3n/Xom/shared"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"strings"
)

type rootHandler struct {
	env      *shared.Env
	telegram shared.TelegramService
	db       *gorm.DB
}

func (h *rootHandler) HandleMessages(update tgbotapi.Update) error {
	message := update.Message
	user, userErr := checkIfUserExists(message.From, h.db)
	if userErr != nil {
		return userErr
	}

	var err error
	if message.IsCommand() {
		command := update.Message.Command()
		switch command {
		case "start":
			err = h.HandleStartCommand(user, message)
		case "help":
			err = h.HandleHelpCommand(user, message)
		}
	}

	if message.Photo != nil || (message.Document != nil && strings.Contains(message.Document.MimeType, "image/")) {
		err = h.HandlePhotoCommand(user, message)
	}
	return err
}

func NewRootHandler(env *shared.Env, telegram shared.TelegramService, db *gorm.DB,
) shared.RootHandlerInterface {
	return &rootHandler{env: env, telegram: telegram, db: db}
}

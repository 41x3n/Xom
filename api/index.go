package api

import (
	"strings"

	interfaces "github.com/41x3n/Xom/interface"
	"github.com/41x3n/Xom/shared"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type rootHandler struct {
	env      *interfaces.Env
	telegram interfaces.TelegramService
	rabbitMQ interfaces.RabbitMQService
	db       *gorm.DB
}

func (h *rootHandler) HandleMessages(update tgbotapi.Update, updateType shared.UpdateType) error {
	var message *tgbotapi.Message
	var callbackQuery *tgbotapi.CallbackQuery
	if updateType == shared.Message {
		message = update.Message

	} else if updateType == shared.Callback {
		message = update.CallbackQuery.Message
		callbackQuery = update.CallbackQuery
	} else {
		return nil
	}

	user, userErr := checkIfUserExists(message.From, h.db)
	if userErr != nil {
		return userErr
	}

	var err error
	if callbackQuery != nil {
		err = h.HandleCallback(user, callbackQuery)
		return err
	}

	if message.IsCommand() {
		command := shared.CommandType(update.Message.Command())
		switch command {
		case shared.StartCommand:
			err = h.HandleStartCommand(user, message)
		case shared.HelpCommand:
			err = h.HandleHelpCommand(user, message)
		}
		return err
	}

	if message.Photo != nil || (message.Document != nil && strings.Contains(message.Document.MimeType, "image/")) {
		err = h.HandlePhoto(user, message)
	}

	if message.Voice != nil || (message.Document != nil && strings.Contains(message.Document.MimeType, "audio/")) {
		err = h.HandleAudio(user, message)
	}

	return err
}

func NewRootHandler(env *interfaces.Env, telegram interfaces.TelegramService,
	rabbitMQ interfaces.RabbitMQService,
	db *gorm.DB,
) interfaces.RootHandlerInterface {
	return &rootHandler{env: env, telegram: telegram, rabbitMQ: rabbitMQ, db: db}
}

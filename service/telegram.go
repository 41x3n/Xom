package service

import (
	"log"

	"github.com/41x3n/Xom/shared"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type telegram struct {
	API *tgbotapi.BotAPI
}

func (t *telegram) PollForUpdates(rootHandler shared.RootHandlerInterface) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.API.GetUpdatesChan(u)

	for update := range updates {
		var updateType shared.UpdateType
		if update.Message != nil {
			updateType = shared.Message
		} else if update.CallbackQuery != nil {
			updateType = shared.Callback
		}

		if updateType != "" {
			if err := rootHandler.HandleMessages(update, updateType); err != nil {
				log.Printf("Error handling update: %v", errors.Wrap(err, "error handling update"))
			}
		}
	}
}

func (t *telegram) GetAPI() *tgbotapi.BotAPI {
	return t.API
}

func NewTelegramService(token string, canDebug bool) (shared.TelegramService, error) {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	botAPI.Debug = canDebug

	return &telegram{
		API: botAPI,
	}, nil
}

package service

import (
	"github.com/41x3n/Xom/shared"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type telegram struct {
	API *tgbotapi.BotAPI
}

func (t *telegram) PollForUpdates(env *shared.Env, rootHandler shared.RootHandlerInterface) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.API.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		// Handle each update. If an error occurs, return it.
		err := rootHandler.HandleMessages(update)
		if err != nil {
			return errors.Wrap(err, "error handling update")
		}
	}

	return nil
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

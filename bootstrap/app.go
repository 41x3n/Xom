package bootstrap

import (
	"github.com/41x3n/Xom/api"
	"github.com/41x3n/Xom/shared"
	"log"
)

type Application struct {
	Env         *shared.Env
	Telegram    shared.TelegramService
	RootHandler shared.RootHandlerInterface
}

type ApplicationInterface interface {
	PollForTelegramUpdates()
}

func (app *Application) PollForTelegramUpdates() {
	err := app.Telegram.PollForUpdates(app.Env, app.RootHandler)
	if err != nil {
		log.Fatalf("error polling for telegram updates: %v", err)
	}
}

func NewApplication() ApplicationInterface {
	app := &Application{}
	app.Env = NewEnv()
	app.Telegram = GetTelegram(app.Env)
	app.RootHandler = api.NewRootHandler(app.Env, app.Telegram)

	return app
}

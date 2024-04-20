package bootstrap

import (
	"github.com/41x3n/Xom/api"
	"github.com/41x3n/Xom/shared"
	"gorm.io/gorm"
)

type Application struct {
	Env         *shared.Env
	Telegram    shared.TelegramService
	RootHandler shared.RootHandlerInterface
	Postgres    *gorm.DB
}

type ApplicationInterface interface {
	PollForTelegramUpdates()
	AutoMigrate()
	CloseDBConnection()
}

func (app *Application) PollForTelegramUpdates() {
	app.Telegram.PollForUpdates(app.RootHandler)
}

func NewApplication() ApplicationInterface {
	app := &Application{}
	app.Env = NewEnv()
	app.Telegram = GetTelegram(app.Env)
	app.Postgres = NewPostgresDatabase(app.Env)
	app.RootHandler = api.NewRootHandler(app.Env, app.Telegram, app.Postgres)

	return app
}

func (app *Application) AutoMigrate() {
	AutoMigrate(app.Postgres)
}

func (app *Application) CloseDBConnection() {
	ClosePostgresDBConnection(app.Postgres)
}

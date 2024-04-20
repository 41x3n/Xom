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
	RabbitMQ    shared.RabbitMQService
}

type ApplicationInterface interface {
	PollForTelegramUpdates()
	AutoMigrate()
	CloseDBConnection()
	CloseRabbitMQ()
	ConsumeMessages()
}

func (app *Application) PollForTelegramUpdates() {
	app.Telegram.PollForUpdates(app.RootHandler)
}

func (app *Application) ConsumeMessages() {
	app.RabbitMQ.ConsumeMessages()
}

func NewApplication() ApplicationInterface {
	app := &Application{}
	app.Env = NewEnv()
	app.Telegram = GetTelegram(app.Env)
	app.RabbitMQ = NewRabbitMQ(app.Env)
	app.Postgres = NewPostgresDatabase(app.Env)
	app.RootHandler = api.NewRootHandler(app.Env, app.Telegram,
		app.RabbitMQ, app.Postgres)

	return app
}

func (app *Application) AutoMigrate() {
	AutoMigrate(app.Postgres)
}

func (app *Application) CloseDBConnection() {
	ClosePostgresDBConnection(app.Postgres)
}

func (app *Application) CloseRabbitMQ() {
	CloseRabbitMQChannel(app.RabbitMQ.GetChannel())
	CloseRabbitMQConnection(app.RabbitMQ.GetConnection())
}

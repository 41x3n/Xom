package bootstrap

import (
	"github.com/41x3n/Xom/api"
	interfaces "github.com/41x3n/Xom/interface"
	service "github.com/41x3n/Xom/service/ffmpeg"
	"gorm.io/gorm"
)

type Application struct {
	Env           *interfaces.Env
	Telegram      interfaces.TelegramService
	RootHandler   interfaces.RootHandlerInterface
	Postgres      *gorm.DB
	RabbitMQ      interfaces.RabbitMQService
	FfmpegHandler interfaces.FFMPEGService
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
	app.Postgres = NewPostgresDatabase(app.Env)
	app.FfmpegHandler = service.NewFFMPEGService(app.Env, app.Telegram,
		app.Postgres)
	app.RabbitMQ = NewRabbitMQ(app.Env, app.FfmpegHandler)
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

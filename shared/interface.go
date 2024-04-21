package shared

import (
	"github.com/41x3n/Xom/core/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Env struct {
	AppEnv           string `mapstructure:"APP_ENV"`
	DSN              string `mapstructure:"DSN"`
	TelegramBotToken string `mapstructure:"TELEGRAM_BOT_TOKEN"`
	RabbitMQURL      string `mapstructure:"RABBITMQ_URL"`
	ContextTimeout   int    `mapstructure:"CONTEXT_TIMEOUT"`
}

type TelegramService interface {
	GetAPI() *tgbotapi.BotAPI
	PollForUpdates(rootHandler RootHandlerInterface)
}

type RabbitMQService interface {
	GetChannel() *amqp.Channel
	GetConnection() *amqp.Connection
	GetQueue() (*amqp.Queue, error)
	PublishMessage(payload RabbitMQPayload) error
	ConsumeMessages()
}

type FFMPEGService interface {
	HandleFiles(payload *RabbitMQPayload) error
}

type RootHandlerInterface interface {
	HandleMessages(update tgbotapi.Update, updateType UpdateType) error
	HandleStartCommand(user *domain.User, message *tgbotapi.Message) error
	HandleHelpCommand(user *domain.User, message *tgbotapi.Message) error
	HandlePhoto(user *domain.User, message *tgbotapi.Message) error
	HandleCallback(user *domain.User, callback *tgbotapi.CallbackQuery) error
}

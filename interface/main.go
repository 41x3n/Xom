package interfaces

import (
	"github.com/41x3n/Xom/core/domain"
	"github.com/41x3n/Xom/shared"
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
	PublishMessage(payload shared.RabbitMQPayload) error
	ConsumeMessages()
}

type FFMPEGService interface {
	HandleFiles(payload *shared.RabbitMQPayload) error
	HandlePhotos(ID int64) error
	ConvertFile(inputPath, outputPath string) error
	ConvertImageToPDF(inputPath, outputPath string) error
	SendFileToUser(media interface{}, outputPath,
		message string) error
	GetInputOutputFilePaths(media interface{}) (string, string, error)
	IsValidFormat(format string) bool
	InformUserAboutError(userTelegramID, messageID int64,
		errorText string) error
}

type RootHandlerInterface interface {
	HandleMessages(update tgbotapi.Update, updateType shared.UpdateType) error
	HandleStartCommand(user *domain.User, message *tgbotapi.Message) error
	HandleHelpCommand(user *domain.User, message *tgbotapi.Message) error
	HandlePhoto(user *domain.User, message *tgbotapi.Message) error
	HandleAudio(user *domain.User, message *tgbotapi.Message) error
	HandleCallback(user *domain.User, callback *tgbotapi.CallbackQuery) error
}

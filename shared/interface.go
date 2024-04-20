package shared

import (
	"github.com/41x3n/Xom/core/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Env struct {
	AppEnv           string `mapstructure:"APP_ENV"`
	TelegramBotToken string `mapstructure:"TELEGRAM_BOT_TOKEN"`
	DSN              string `mapstructure:"DSN"`
	ContextTimeout   int    `mapstructure:"CONTEXT_TIMEOUT"`
}

type TelegramService interface {
	GetAPI() *tgbotapi.BotAPI
	PollForUpdates(rootHandler RootHandlerInterface)
}

type RootHandlerInterface interface {
	HandleMessages(update tgbotapi.Update, updateType UpdateType) error
	HandleStartCommand(user *domain.User, message *tgbotapi.Message) error
	HandleHelpCommand(user *domain.User, message *tgbotapi.Message) error
	HandlePhoto(user *domain.User, message *tgbotapi.Message) error
	HandleCallback(user *domain.User, callback *tgbotapi.CallbackQuery) error
}

package shared

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Env struct {
	AppEnv           string `mapstructure:"APP_ENV"`
	TelegramBotToken string `mapstructure:"TELEGRAM_BOT_TOKEN"`
}

type TelegramService interface {
	GetAPI() *tgbotapi.BotAPI
	PollForUpdates(env *Env, rootHandler RootHandlerInterface) error
}

type RootHandlerInterface interface {
	HandleMessages(update tgbotapi.Update) error
}

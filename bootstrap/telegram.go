package bootstrap

import (
	"log"

	interfaces "github.com/41x3n/Xom/interface"
	"github.com/41x3n/Xom/service"
)

func GetTelegram(env *interfaces.Env) interfaces.TelegramService {
	token := env.TelegramBotToken
	canDebug := env.AppEnv == "development"

	telegram, err := service.NewTelegramService(token, canDebug)

	if err != nil {
		log.Fatal("Can't connect to Telegram: ", err)
	}

	API := telegram.GetAPI()

	log.Println("Connection to Telegram established.")
	log.Println("Authorized on account ", API.Self.UserName)

	return telegram
}

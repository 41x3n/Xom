package bootstrap

import (
	"github.com/41x3n/Xom/shared"
	"log"

	"github.com/spf13/viper"
)

func NewEnv() *shared.Env {
	env := shared.Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	return &env
}

package bootstrap

import (
	"log"

	interfaces "github.com/41x3n/Xom/interface"
	"github.com/spf13/viper"
)

func NewEnv() *interfaces.Env {
	env := interfaces.Env{}
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

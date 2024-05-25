package main

import (
	"log"
	"os"

	"github.com/41x3n/Xom/bootstrap"
)

func main() {
	// Customize the log output format
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Llongfile)

	app := bootstrap.NewApplication()

	defer app.CloseDBConnection()
	defer app.CloseRabbitMQ()

	app.AutoMigrate()

	go app.ConsumeMessages()

	// Start polling Telegram in a goroutine
	go app.PollForTelegramUpdates()

	// Wait indefinitely
	select {}
}

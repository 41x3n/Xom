package bootstrap

import (
	"github.com/41x3n/Xom/service"
	"github.com/41x3n/Xom/shared"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func CreateRabbitMQConnection(env *shared.Env) *amqp.Connection {
	url := env.RabbitMQURL

	conn, err := amqp.Dial(url)
	shared.FailOnError(err, "Failed to connect to RabbitMQ")

	log.Println("Connection to RabbitMQ established.")
	return conn
}

func CloseRabbitMQConnection(conn *amqp.Connection) {
	conn.Close()
}

func CreateChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	shared.FailOnError(err, "Failed to open a channel")
	return ch
}

func CloseRabbitMQChannel(ch *amqp.Channel) {
	ch.Close()
}

func NewRabbitMQ(env *shared.Env, ffmpegHandler shared.FFMPEGService) shared.
	RabbitMQService {
	conn := CreateRabbitMQConnection(env)
	ch := CreateChannel(conn)

	rabbitMQ, err := service.NewRabbitMQService(conn, ch, env, ffmpegHandler)

	shared.FailOnError(err, "Failed to create RabbitMQ service")

	return rabbitMQ
}

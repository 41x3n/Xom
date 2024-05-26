package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	interfaces "github.com/41x3n/Xom/interface"
	"github.com/41x3n/Xom/shared"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	q             amqp.Queue
	queueDeclared bool
)

const FileQueue = "file_queue"

type rabbitMQ struct {
	Conn          *amqp.Connection
	Ch            *amqp.Channel
	env           *interfaces.Env
	FfmpegHandler interfaces.FFMPEGService
}

func (r *rabbitMQ) GetChannel() *amqp.Channel {
	return r.Ch
}

func (r *rabbitMQ) GetConnection() *amqp.Connection {
	return r.Conn
}

func (r *rabbitMQ) GetQueue() (*amqp.Queue, error) {
	if !queueDeclared {
		ch := r.Ch

		var err error
		q, err = ch.QueueDeclare(
			FileQueue,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to declare a queue: %v", err)
		}
		queueDeclared = true
	}

	return &q, nil
}

func (r *rabbitMQ) PublishMessage(payload shared.RabbitMQPayload) error {
	ch := r.Ch

	q, err := r.GetQueue()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Serialize the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to serialize payload: %v", err)
	}

	// Retry mechanism
	for i := 0; i < 3; i++ {
		err = ch.PublishWithContext(
			ctx,
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        jsonPayload,
			},
		)

		if err == nil {
			return nil
		}

		// Sleep before retrying
		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("failed to publish a message: %v", err)
}

func (r *rabbitMQ) ConsumeMessages() {
	ch := r.Ch

	q, err := r.GetQueue()
	if err != nil {
		log.Fatalf("failed to get queue: %v", err)
		return
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("failed to register a consumer: %v", err)
	}

	sem := make(chan struct{}, 10) // Create a semaphore channel with a capacity of 10

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
		go r.processMessage(d, sem)
	}
}

func (r *rabbitMQ) processMessage(delivery amqp.Delivery, semaphore chan struct{}) {
	semaphore <- struct{}{} // Send to the semaphore channel

	payload := &shared.RabbitMQPayload{}
	if errMarshal := json.Unmarshal(delivery.Body, payload); errMarshal != nil {
		log.Printf("Error decoding JSON: %v", errMarshal)
		return
	}

	if err := r.FfmpegHandler.HandleFiles(payload); err != nil {
		log.Printf("Error handling command: %v", err)
		r.ackMessage(delivery, false)
		return
	}

	r.ackMessage(delivery, true)
	<-semaphore
}

func (r *rabbitMQ) ackMessage(delivery amqp.Delivery, ackType bool) {
	if err := delivery.Ack(ackType); err != nil {
		log.Printf("Error acknowledging message: %v", err)
	}
}

func NewRabbitMQService(conn *amqp.Connection, ch *amqp.Channel,
	env *interfaces.Env, ffmpegHandler interfaces.FFMPEGService) (interfaces.RabbitMQService,
	error) {
	return &rabbitMQ{
		Conn:          conn,
		Ch:            ch,
		env:           env,
		FfmpegHandler: ffmpegHandler,
	}, nil
}

package helper

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type IRabbitMQ interface {
	Publisher(name string, data []byte, exchange string, exchangeType string) error
}

type rabbitmq struct {
	channel *amqp091.Channel
}

func IntializeRabbitMQ() IRabbitMQ {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%s",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASS"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
	)

	connection, err := amqp091.Dial(dsn)
	if err != nil {
		log.Println("Error dial to AMQP")
	}

	defer connection.Close()
	channel, err := connection.Channel()
	if err != nil {
		log.Println("Error while create channel")
	}

	defer channel.Close()

	return &rabbitmq{channel: channel}
}

func (r *rabbitmq) Publisher(name string, data []byte, exchange string, exchangeType string) error {
	log.Println("[*] Declaring Exchange...")
	err := r.channel.ExchangeDeclare(exchange, exchangeType, false, false, false, false, nil)
	if err != nil {
		log.Println("[x] Error declaring exchange")
	}

	log.Println("[*] Declaring Queue...")
	queue, err := r.channel.QueueDeclare(name, true, false, false, false , nil )
	if err != nil {
		log.Println("[x] Error declaring queue")

	}

	log.Println("[*] Queue Binding...")
	err_binding := r.channel.QueueBind(queue.Name, "info", exchange, false, nil)
	if err_binding != nil {
		log.Println("[x] Error binding queue")
	}

	log.Println("[*] Publish With Context...")
	err_pub := r.channel.PublishWithContext(context.Background(), exchange, queue.Name, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body: data,
		Timestamp: time.Now(),
	})

	if err_pub != nil {
		log.Println("[x] Error publishing queue")
	}

	log.Println("[*] Done Publishing")
	return nil
}
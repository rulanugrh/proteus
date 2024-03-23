package pkg

import (
	"context"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/rulanugrh/tokoku/product/internal/config"
	"github.com/rulanugrh/tokoku/product/internal/entity/web"
)

type RabbitMQInterface interface {
	Publish(name string, data []byte, exhange string, exchangeType string, userID string) error
}

type rabbit struct {
	channel config.RabbitMQ
}

func RabbitMQ(channel config.RabbitMQ) RabbitMQInterface {
	return &rabbit{channel: channel}
}

func(r *rabbit) Publish(name string, data []byte, exhange string, exchangeType string, userID string) error {
	log.Println("[*] Declaring Exchange...")
	err := r.channel.Channel.ExchangeDeclare(exhange, exchangeType, false, false, false, false, nil)
	if err != nil {
		return web.InternalServerError(err.Error())
	}

	log.Println("[*] Declaring Queue...")
	queue, err := r.channel.Channel.QueueDeclare(name, true, false, false, false, nil)
	if err != nil {
		return web.InternalServerError(err.Error())
	}

	log.Println("[*] Queue Binding...")
	if err = r.channel.Channel.QueueBind(queue.Name, "info", exhange, false, nil); err != nil {
		return web.BadRequest(err.Error())
	}

	log.Println("[*] Publisher with context...")

	err_pub := r.channel.Channel.PublishWithContext(context.Background(), exhange, queue.Name, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body: data,
		Timestamp: time.Now(),
		UserId: userID,
	})

	if err_pub != nil {
		return web.InternalServerError(err_pub.Error())
	}

	log.Println("[*] Publisher Success")
	return nil
}
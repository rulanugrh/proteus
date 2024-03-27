package pkg

import (
	"context"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/rulanugrh/order/internal/config"
	"github.com/rulanugrh/order/internal/util/constant"
)

type RabbitMQInterface interface {
	Publisher(name string, data []byte, exchange string, exchangeType string, username string) error
}

type rabbit struct {
	client *config.RabbitMQ
}

func RabbitMQ(client *config.RabbitMQ) RabbitMQInterface {
	return &rabbit{client: client}
}


func (r *rabbit) Publisher(name string, data []byte, exchange string, exchangeType string, username string) error {
	log.Println("[*] Declaring Exchange...")
	err := r.client.Channel.ExchangeDeclare(exchange, exchangeType, false, false, false, false, nil)
	if err != nil {
		return constant.InternalServerError("error exchange declare", err)
	}

	log.Println("[*] Declaring Queue...")
	queue, err_queue := r.client.Channel.QueueDeclare(name, true, false, false, false, nil)
	if err_queue != nil {
		return constant.InternalServerError("error declaring queue", err_queue)
	}

	log.Println("[*] Queue Binding...")
	err_binding := r.client.Channel.QueueBind(queue.Name, "info", exchange, false, nil)
	if err_binding != nil {
		return constant.InternalServerError("error binding queue", err_binding)
	}

	log.Println("[*] Publisher with context ...")
	err_pub := r.client.Channel.PublishWithContext(context.Background(), exchange, queue.Name, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body: data,
		Timestamp: time.Now(),
		UserId: username,
	})

	if err_pub != nil {
		return constant.InternalServerError("error publisher data", err_pub)
	}

	log.Println("[*] Publisher Success")
	return nil

}

func (r *rabbit) CatchProduct() {

}
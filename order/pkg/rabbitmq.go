package pkg

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/rulanugrh/order/internal/config"
	"github.com/rulanugrh/order/internal/entity"
	"github.com/rulanugrh/order/internal/repository"
	"github.com/rulanugrh/order/internal/util/constant"
)

type RabbitMQInterface interface {
	Publisher(name string, data []byte, exchange string, exchangeType string, username string) error
}

type rabbit struct {
	client *config.RabbitMQ
	db     repository.ProductInterface
}

func RabbitMQ(client *config.RabbitMQ, db repository.ProductInterface) RabbitMQInterface {
	return &rabbit{
		client: client,
		db:     db,
	}
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
		Body:        data,
		Timestamp:   time.Now(),
		UserId:      username,
	})

	if err_pub != nil {
		return constant.InternalServerError("error publisher data", err_pub)
	}

	log.Println("[*] Publisher Success")
	return nil

}

func (r *rabbit) CatchProduct() error {
	log.Println("[*] Declare Queue for Catch Created Product")

	queue, err_queue := r.client.Channel.QueueDeclare("product-create", true, false, false, false, nil)
	if err_queue != nil {
		return constant.InternalServerError("error declaring queue for catch create product", err_queue)
	}

	log.Println("[*] Consuming Product Create ...")
	message, err_message := r.client.Channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err_message != nil {
		return constant.InternalServerError("error consume for product-create", err_message)
	}

	var response chan struct{}
	go func() {
		for msg := range message {
			log.Println("[*] Success Receive Message")
			var payload entity.Product
			err := json.Unmarshal(msg.Body, &payload)
			if err != nil {
				log.Printf("error marshaling response: %s", err.Error())
			}

			err_created := r.db.Create(payload)
			if err_created != nil {
				log.Printf("error create, because: %s", err_created.Error())
			}

		}
	}()

	log.Println("[*] Success Consume Product Created")
	<-response
	return nil
}

func (r *rabbit) UpdateProduct() error {
	log.Println("[*] Declaring Queue for Update Product")
	queue, err_queue := r.client.Channel.QueueDeclare("product-update", true, false, false, false, nil)
	if err_queue != nil {
		return constant.InternalServerError("error, cannot declare queue", err_queue)
	}

	log.Println("[*] Start Consume ...")
	message, err_message := r.client.Channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err_message != nil {
		return constant.InternalServerError("sorry cannot consume this queue", err_message)
	}

	var response chan struct{}
	go func() {
		for msg := range message {
			log.Println("[*] Receiving Message")
			var payload entity.Product
			err := json.Unmarshal(msg.Body, &payload)
			if err != nil {
				log.Printf("error marshaling response: %s", err.Error())
			}

			err_update := r.db.Update(payload.ID, payload)
			if err_update != nil {
				log.Printf("error create, because: %s", err_update.Error())
			}

		}
	}()

	log.Printf("[*] Success Consume Message ( Update Product )")
	<-response
	return nil
}

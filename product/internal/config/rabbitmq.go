package config

import (
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Channel *amqp091.Channel
	conf *App
}

func InitRabbit(conf *App) *RabbitMQ{
	return &RabbitMQ{conf: conf}
}
func (r *RabbitMQ) InitRabbit() {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%s", 
		r.conf.RabbitMQ.User,
		r.conf.RabbitMQ.Pass,
		r.conf.RabbitMQ.Host, 
		r.conf.RabbitMQ.Port,
	)

	connection, err := amqp091.Dial(dsn)
	if err != nil {
		log.Printf("error, cant connect to amqp :%s", err.Error())
	}

	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		log.Printf("error, cant connect to amqp :%s", err.Error())
	}

	defer channel.Close()

	r.Channel = channel
}
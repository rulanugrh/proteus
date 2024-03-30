package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rulanugrh/webhook/handler"
	"github.com/rulanugrh/webhook/helper"
	"github.com/rulanugrh/webhook/routes"
)

func main() {
	rabbitmq := helper.IntializeRabbitMQ()
	dsn := fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT"))
	webhook := handler.NewWebhookHandler(rabbitmq)

	err := routes.HandlerRoute(webhook, dsn)
	if err != nil {
		log.Println("[*] Error: ", err.Error())
	}

	log.Printf("Webhook running at :%s", dsn)
}
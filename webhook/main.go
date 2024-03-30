package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rulanugrh/webhook/handler"
	"github.com/rulanugrh/webhook/routes"
)

func main() {
	dsn := fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT"))
	webhook := handler.NewWebhookHandler()

	err := routes.HandlerRoute(webhook, dsn)
	if err != nil {
		log.Println("[*] Error: ", err.Error())
	}

	log.Printf("Webhook running at :%s", dsn)
}
package routes

import (
	"net/http"

	"github.com/rulanugrh/webhook/handler"
)

func HandlerRoute(webhook *handler.WebhookHandler, addr string) error {
	serv := http.NewServeMux()
	serv.HandleFunc("/payment_success", webhook.PaymentSucess)
	serv.HandleFunc("/payment_failed", webhook.PaymentFailed)
	serv.HandleFunc("/payment_waiting", webhook.AwaitPaymentCapture)
	serv.HandleFunc("/payment_failed", webhook.PaymentFailed)

	server := http.Server{
		Addr: addr,
		Handler: serv,
	}

	return server.ListenAndServe()
}
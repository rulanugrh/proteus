package handler

import "net/http"

type WebhookInterface interface {
	PaymentSuccess(w http.ResponseWriter, r *http.Request)
}

type webhook struct {}

func NewWebhookHandler() WebhookInterface {
	return &webhook{}
}

func (wb *webhook) PaymentSuccess(w http.ResponseWriter, r *http.Request) {

}
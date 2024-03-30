package handler

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/rulanugrh/webhook/helper"
)

type WebhookHandler struct {
	rabbit helper.IRabbitMQ
}

func NewWebhookHandler(rabbit helper.IRabbitMQ) *WebhookHandler {
	return &WebhookHandler{rabbit: rabbit}
}

func (wh *WebhookHandler) PaymentSucess(w http.ResponseWriter, r *http.Request) {
	
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Failure read request body")
		w.WriteHeader(400)
		w.Write(helper.Marshal(helper.BadRequest(err.Error(), "bad request")))
		return 
	}
	
	log.Println("[*] Webhook Payment Success")
	log.Println("[*] Response")
	log.Println(string(body))

	response := helper.Marshal(helper.Success("payment success", "ok"))
	err = wh.rabbit.Publisher("notif-webhook", response, "webhook", "topic")
	if err != nil {
		log.Println("[X] Failure Publishing")
	}

	w.WriteHeader(200)
	w.Write(response)
	return
}

func (wh *WebhookHandler) AwaitPaymentCapture(w http.ResponseWriter, r *http.Request) {
	
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Failure read request body")
		w.WriteHeader(400)
		w.Write(helper.Marshal(helper.BadRequest(err.Error(), "bad request")))
		return 
	}
	
	log.Println("[*] Await Payment Capture")
	log.Println("[*] Response")
	log.Println(string(body))

	response := helper.Marshal(helper.Success("await payment capture", "ok"))
	err = wh.rabbit.Publisher("notif-webhook", response, "webhook", "topic")
	if err != nil {
		log.Println("[X] Failure Publishing")
	}

	w.WriteHeader(200)
	w.Write(response)
	return
}

func (wh *WebhookHandler) PaymentPending(w http.ResponseWriter, r *http.Request) {
	
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Failure read request body")
		w.WriteHeader(400)
		w.Write(helper.Marshal(helper.BadRequest(err.Error(), "bad request")))
		return 
	}
	
	log.Println("[*] Webhook Payment Failed")
	log.Println("[*] Response")
	log.Println(string(body))

	response := helper.Marshal(helper.Success("payment pending", "ok"))
	err = wh.rabbit.Publisher("notif-webhook", response, "webhook", "topic")
	if err != nil {
		log.Println("[X] Failure Publishing")
	}

	w.WriteHeader(200)
	w.Write(response)
	return
}

func (wh *WebhookHandler) PaymentFailed(w http.ResponseWriter, r *http.Request) {
	
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Failure read request body")
		w.WriteHeader(400)
		w.Write(helper.Marshal(helper.BadRequest(err.Error(), "bad request")))
		return 
	}
	
	log.Println("[*] Webhook Payment Failed")
	log.Println("[*] Response")
	log.Println(string(body))

	response := helper.Marshal(helper.Success("payment failed", "ok"))
	err = wh.rabbit.Publisher("notif-webhook", response, "webhook", "topic")
	if err != nil {
		log.Println("[X] Failure Publishing")
	}
	
	w.WriteHeader(200)
	w.Write(response)
	return
}
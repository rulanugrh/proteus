package handler

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/rulanugrh/webhook/helper"
)

type WebhookHandler struct{}

func NewWebhookHandler() *WebhookHandler {
	return &WebhookHandler{}
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

	w.WriteHeader(200)
	w.Write(helper.Marshal(helper.Success("payment success", "ok")))
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

	w.WriteHeader(200)
	w.Write(helper.Marshal(helper.Success("await payment capture", "ok")))
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

	w.WriteHeader(200)
	w.Write(helper.Marshal(helper.Success("payment pending", "ok")))
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

	w.WriteHeader(200)
	w.Write(helper.Marshal(helper.Success("payment failed", "ok")))
	return
}
package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/rulanugrh/webhook/helper"
)

type Response struct {
	Event      string `json:"event"`
	BusinessID string `json:"business_id"`
	Data       struct {
		ID          string `json:"id"`
		Amount      string `json:"amout"`
		Country     string `json:"ID"`
		Currency    string `json:"IDR"`
		ReferenceID string `json:"reference_id"`
		Status      string `json:"status"`
	}
}

type WebhookHandler struct {
	rabbit helper.IRabbitMQ
}

func NewWebhookHandler(rabbit helper.IRabbitMQ) *WebhookHandler {
	return &WebhookHandler{rabbit: rabbit}
}

func (wh *WebhookHandler) FVAPaid(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("[x] Error read body")
		w.WriteHeader(400)
		w.Write(helper.Marshal(helper.BadRequest(err.Error(), "bad request")))
		return
	}
	log.Println("[*] Webhook FVA Paid")
	log.Println("[*] Response")
	log.Println(string(body))

	response := helper.Marshal(helper.Success("fva paid", "ok", nil))
	w.WriteHeader(200)
	w.Write(response)
	return
}

func (wh *WebhookHandler) FVACreated(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("[x] Error read body")
		w.WriteHeader(400)
		w.Write(helper.Marshal(helper.BadRequest(err.Error(), "bad request")))
		return
	}
	log.Println("[*] Webhook FVA Created")
	log.Println("[*] Response")
	log.Println(string(body))

	response := helper.Marshal(helper.Success("fva created", "ok", nil))
	w.WriteHeader(200)
	w.Write(response)
	return
}

func (wh *WebhookHandler) CaptureSuccess(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("[x] Error read body")
		w.WriteHeader(400)
		w.Write(helper.Marshal(helper.BadRequest(err.Error(), "bad request")))
		return
	}
	log.Println("[*] Webhook Capture Success")
	log.Println("[*] Response")
	log.Println(string(body))

	response := helper.Marshal(helper.Success("capture success", "ok", nil))
	w.WriteHeader(200)
	w.Write(response)
	return
}

func (wh *WebhookHandler) CaptureFailed(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("[x] Error read body")
		w.WriteHeader(400)
		w.Write(helper.Marshal(helper.BadRequest(err.Error(), "bad request")))
		return
	}
	log.Println("[*] Webhook Capture Failed")
	log.Println("[*] Response")
	log.Println(string(body))

	response := helper.Marshal(helper.Success("capture failed", "ok", nil))
	w.WriteHeader(200)
	w.Write(response)
	return
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

	read_data := json.Unmarshal(body, &Response{})
	response := helper.Marshal(helper.Success("payment success", "ok", read_data))
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

	read_data := json.Unmarshal(body, &Response{})
	response := helper.Marshal(helper.Success("await payment capture", "ok", read_data))
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

	read_data := json.Unmarshal(body, &Response{})
	response := helper.Marshal(helper.Success("payment pending", "ok", read_data))
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

	read_data := json.Unmarshal(body, &Response{})
	response := helper.Marshal(helper.Success("payment failed", "ok", read_data))
	err = wh.rabbit.Publisher("notif-webhook", response, "webhook", "topic")
	if err != nil {
		log.Println("[X] Failure Publishing")
	}

	w.WriteHeader(200)
	w.Write(response)
	return
}

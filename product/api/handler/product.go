package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/rulanugrh/tokoku/product/internal/entity/domain"
	"github.com/rulanugrh/tokoku/product/internal/entity/web"
	"github.com/rulanugrh/tokoku/product/internal/middleware"
	"github.com/rulanugrh/tokoku/product/internal/service"
	"github.com/rulanugrh/tokoku/product/pkg"
)

type ProductInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
	FindID(w http.ResponseWriter, r *http.Request)
	FindAll(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
}

type product struct {
	service service.ProductInterface
	rabbitmq pkg.RabbitMQInterface
}

func ProductHandler(service service.ProductInterface, rabbitmq pkg.RabbitMQInterface) ProductInterface {
	return &product{service: service, rabbitmq: rabbitmq}
}

func(p *product) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.Product

	token := r.Header.Get("Authorization")

	claim, err := middleware.CheckToken(token)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	err = middleware.ValidateRole(token)
	if err != nil {
		response, err := json.Marshal(web.Forbidden(err.Error()))
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(403)
		w.Write(response)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response, err := json.Marshal(web.InternalServerError("sorry cannot read body request"))
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(500)
		w.Write(response)
		return
	}

	data, err := p.service.Create(req)
	if err != nil {
		response, err := json.Marshal(web.BadRequest(err.Error()))
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(400)
		w.Write(response)
		return
	}

	marshalling, _ := json.Marshal(data)

	err = p.rabbitmq.Publish("product-create", marshalling, "product-exchange", "topic", claim.Username)
	if err != nil {
		response, err := json.Marshal(web.BadRequest(err.Error()))
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(400)
		w.Write(response)
		return
	}

	response, err := json.Marshal(web.Created(data, "success create product"))
	if err != nil {
		w.WriteHeader(500)
		return
	}


	w.WriteHeader(201)
	w.Write(response)
	return
}

func(p *product) FindID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/product/find/"))
	if err != nil {
		w.WriteHeader(500)
		return
	}

	data, err := p.service.FindID(uint(id))
	if err != nil {
		response, err := json.Marshal(web.BadRequest(err.Error()))
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(400)
		w.Write(response)
		return
	}

	response, err := json.Marshal(web.Success(data, "success get product"))
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	w.Write(response)
	return
}

func(p *product) FindAll(w http.ResponseWriter, r *http.Request) {
	data, err := p.service.FindAll()
	if err != nil {
		response, err := json.Marshal(web.BadRequest(err.Error()))
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(400)
		w.Write(response)
		return
	}

	response, err := json.Marshal(web.Success(data, "success get product"))
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	w.Write(response)
	return
}

func(p *product) Update(w http.ResponseWriter, r *http.Request) {
	var req domain.Product
	token := r.Header.Get("Authorization")

	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/product/update/"))
	if err != nil {
		w.WriteHeader(500)
		return
	}

	err = middleware.ValidateRole(token)
	if err != nil {
		response, err := json.Marshal(web.Forbidden(err.Error()))
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(403)
		w.Write(response)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response, err := json.Marshal(web.InternalServerError("sorry cannot read body request"))
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(500)
		w.Write(response)
		return
	}
	data, err := p.service.Update(uint(id), req)
	if err != nil {
		response, err := json.Marshal(web.BadRequest(err.Error()))
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(400)
		w.Write(response)
		return
	}

	response, err := json.Marshal(web.Success(data, "success update product"))
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	w.Write(response)
	return
}

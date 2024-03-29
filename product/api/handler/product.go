package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
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
	metric *pkg.Metric
}

func ProductHandler(service service.ProductInterface, rabbitmq pkg.RabbitMQInterface, metric *pkg.Metric) ProductInterface {
	return &product{service: service, rabbitmq: rabbitmq, metric: metric}
}

func(p *product) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.Product

	token := r.Header.Get("Authorization")

	claim, err := middleware.CheckToken(token)
	if err != nil {
		p.metric.Histogram.With(prometheus.Labels{"code": "401", "method": "POST", "type": "create", "service": "product"}).Observe(time.Since(time.Now()).Seconds())
		
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
		p.metric.Histogram.With(prometheus.Labels{"code": "403", "method": "POST", "type": "create", "service": "product"}).Observe(time.Since(time.Now()).Seconds())

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
		p.metric.Histogram.With(prometheus.Labels{"code": "500", "method": "POST", "type": "create", "service": "product"}).Observe(time.Since(time.Now()).Seconds())

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
		p.metric.Histogram.With(prometheus.Labels{"code": "400", "method": "POST", "type": "create", "service": "product"}).Observe(time.Since(time.Now()).Seconds())

		w.WriteHeader(400)
		w.Write(response)
		return
	}

	marshalling, _ := json.Marshal(data)

	err = p.rabbitmq.Publish("product-create", marshalling, "product", "topic", claim.Username)
	if err != nil {
		response, err := json.Marshal(web.BadRequest(err.Error()))
		if err != nil {
			w.WriteHeader(500)
			return
		}

		p.metric.Histogram.With(prometheus.Labels{"code": "400", "method": "POST", "type": "create", "service": "product"}).Observe(time.Since(time.Now()).Seconds())

		w.WriteHeader(400)
		w.Write(response)
		return
	}

	response, err := json.Marshal(web.Created(data, "success create product"))
	if err != nil {
		w.WriteHeader(500)
		return
	}

	p.metric.Histogram.With(prometheus.Labels{"code": "201", "method": "POST", "type": "create", "service": "product"}).Observe(time.Since(time.Now()).Seconds())
	p.metric.Counter.With(prometheus.Labels{"type": "create", "service": "product"}).Inc()
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

		p.metric.Histogram.With(prometheus.Labels{"code": "400", "method": "GET", "type": "findID", "service": "product"}).Observe(time.Since(time.Now()).Seconds())

		w.WriteHeader(400)
		w.Write(response)
		return
	}

	response, err := json.Marshal(web.Success(data, "success get product"))
	if err != nil {
		w.WriteHeader(500)
		return
	}

	p.metric.Histogram.With(prometheus.Labels{"code": "200", "method": "GET", "type": "findID", "service": "product"}).Observe(time.Since(time.Now()).Seconds())
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
		p.metric.Histogram.With(prometheus.Labels{"code": "400", "method": "GET", "type": "findAll", "service": "product"}).Observe(time.Since(time.Now()).Seconds())

		w.WriteHeader(400)
		w.Write(response)
		return
	}

	response, err := json.Marshal(web.Success(data, "success get product"))
	if err != nil {
		w.WriteHeader(500)
		return
	}

	p.metric.Histogram.With(prometheus.Labels{"code": "200", "method": "GET", "type": "findAll", "service": "product"}).Observe(time.Since(time.Now()).Seconds())
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
		p.metric.Histogram.With(prometheus.Labels{"code": "403", "method": "PUT", "type": "update", "service": "product"}).Observe(time.Since(time.Now()).Seconds())

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
		p.metric.Histogram.With(prometheus.Labels{"code": "500", "method": "PUT", "type": "update", "service": "product"}).Observe(time.Since(time.Now()).Seconds())

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
		p.metric.Histogram.With(prometheus.Labels{"code": "400", "method": "PUT", "type": "update", "service": "product"}).Observe(time.Since(time.Now()).Seconds())

		w.WriteHeader(400)
		w.Write(response)
		return
	}

	response, err := json.Marshal(web.Success(data, "success update product"))
	if err != nil {
		w.WriteHeader(500)
		return
	}

	p.metric.Histogram.With(prometheus.Labels{"code": "200", "method": "PUT", "type": "update", "service": "product"}).Observe(time.Since(time.Now()).Seconds())
	p.metric.Counter.With(prometheus.Labels{"type": "update", "service": "product"}).Inc()
	w.WriteHeader(200)
	w.Write(response)
	return
}

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
	proto "google.golang.org/protobuf/proto"
)

type ProductInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
	FindID(w http.ResponseWriter, r *http.Request)
	FindAll(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
}

type product struct {
	service  service.ProductInterface
	rabbitmq pkg.RabbitMQInterface
	metric   *pkg.Metric
	log      pkg.ILogrus
}

func ProductHandler(service service.ProductInterface, rabbitmq pkg.RabbitMQInterface, metric *pkg.Metric) ProductInterface {
	return &product{service: service, rabbitmq: rabbitmq, metric: metric, log: pkg.Logrus()}
}

func (p *product) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.Product

	token := r.Header.Get("Authorization")

	claim, err := middleware.CheckToken(token)
	if err != nil {

		p.log.Record("/api/product/create", 401, "POST").Warn("not have token")
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

		p.log.Record("/api/product/create", 403, "POST").Error("forbidden for this user")
		p.metric.Histogram.With(prometheus.Labels{"code": "403", "method": "POST", "type": "create", "service": "product"}).Observe(time.Since(time.Now()).Seconds())

		w.WriteHeader(403)
		w.Write(response)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	data, err := p.service.Create(req)
	if err != nil {
		response, err := json.Marshal(web.BadRequest(err.Error()))
		if err != nil {
			w.WriteHeader(500)
			return
		}

		p.log.Record("/api/product/create", 400, "POST").Error("bad request for create data")
		p.metric.Histogram.With(prometheus.Labels{"code": "400", "method": "POST", "type": "create", "service": "product"}).Observe(time.Since(time.Now()).Seconds())

		w.WriteHeader(400)
		w.Write(response)
		return
	}

	events := web.Product{
		ID: data.ID,
		Name: data.Name,
		Description: data.Description,
		Price: data.Price,
		Available: data.Available,
		Reserved: data.Reserved,
		Category: data.Category,
	}

	marshalling, _ := proto.Marshal(&events)

	err = p.rabbitmq.Publish("product-create", marshalling, "product.info", "topic", claim.Username)
	if err != nil {
		response, err := json.Marshal(web.BadRequest(err.Error()))
		if err != nil {
			w.WriteHeader(500)
			return
		}

		p.log.Record("/api/product/create", 400, "POST").Error("error publish to message broker")
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
}

func (p *product) FindID(w http.ResponseWriter, r *http.Request) {
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

		p.log.Record("/api/product/find/"+strconv.Itoa(id), 400, "GET").Error("cannot get product by this id")
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
}

func (p *product) FindAll(w http.ResponseWriter, r *http.Request) {
	data, err := p.service.FindAll()
	if err != nil {
		response, err := json.Marshal(web.BadRequest(err.Error()))
		if err != nil {
			w.WriteHeader(500)
			return
		}

		p.log.Record("/api/product/get", 400, "GET").Error("error get all product")
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
}

func (p *product) Update(w http.ResponseWriter, r *http.Request) {
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

		p.log.Record("/api/product/update"+strconv.Itoa(id), 403, "PUT").Error("forbidden for this user")
		p.metric.Histogram.With(prometheus.Labels{"code": "403", "method": "PUT", "type": "update", "service": "product"}).Observe(time.Since(time.Now()).Seconds())

		w.WriteHeader(403)
		w.Write(response)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	data, err := p.service.Update(uint(id), req)
	if err != nil {
		response, err := json.Marshal(web.BadRequest(err.Error()))
		if err != nil {
			w.WriteHeader(500)
			return
		}

		p.log.Record("/api/product/update"+strconv.Itoa(id), 400, "PUT").Error("bad request for update data with this id")
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
}

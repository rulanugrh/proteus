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

type CategoryInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
	FindID(w http.ResponseWriter, r *http.Request)
	FindAll(w http.ResponseWriter, r *http.Request)
}

type category struct {
	service service.CategoryInterface
	metric *pkg.Metric
}

func CategoryHandler(service service.CategoryInterface, metric *pkg.Metric) CategoryInterface {
	return &category{service: service, metric: metric}
}

func(c *category) Create(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	err := middleware.ValidateRole(token)
	if err != nil {
		response, err := json.Marshal(web.Forbidden(err.Error()))
		if err != nil {
			w.WriteHeader(500)
			return
		}
		
		c.metric.Histogram.With(prometheus.Labels{"code": "403", "method": "POST", "type": "create", "service": "category"}).Observe(time.Since(time.Now()).Seconds())
		w.WriteHeader(403)
		w.Write(response)
		return
	}

	var req domain.Category
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	data, err := c.service.Create(req)
	if err != nil {
		c.metric.Histogram.With(prometheus.Labels{"code": "400", "method": "POST", "type": "create", "service": "category"}).Observe(time.Since(time.Now()).Seconds())
		response, err := json.Marshal(web.BadRequest(err.Error()))
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(400)
		w.Write(response)
		return
	}

	response, err := json.Marshal(web.Created(data, "success create category"))
	if err != nil {
		w.WriteHeader(500)
		return
	}

	c.metric.Counter.With(prometheus.Labels{"type": "create", "service": "category"}).Inc()
	c.metric.Histogram.With(prometheus.Labels{"code": "201", "method": "POST", "type": "create", "service": "category"}).Observe(time.Since(time.Now()).Seconds())
	
	w.WriteHeader(201)
	w.Write(response)
	return
}

func(c *category) FindID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/category/find/"))
	if err != nil {
		w.WriteHeader(500)
		return
	}

	data, err := c.service.FindID(uint(id))
	if err != nil {
		response, err := json.Marshal(web.BadRequest(err.Error()))
		if err != nil {
			w.WriteHeader(500)
			return
		}

		c.metric.Histogram.With(prometheus.Labels{"code": "400", "method": "GET", "type": "findID", "service": "category"}).Observe(time.Since(time.Now()).Seconds())

		w.WriteHeader(400)
		w.Write(response)
		return
	}

	response, err := json.Marshal(web.Success(data, "category found"))
	if err != nil {
		w.WriteHeader(500)
		return
	}

	c.metric.Histogram.With(prometheus.Labels{"code": "200", "method": "GET", "type": "findID", "service": "category"}).Observe(time.Since(time.Now()).Seconds())
	w.WriteHeader(200)
	w.Write(response)
	return
}

func(c *category) FindAll(w http.ResponseWriter, r *http.Request) {
	data, err := c.service.FindAll()
	if err != nil {
		response, err := json.Marshal(web.BadRequest(err.Error()))
		if err != nil {
			w.WriteHeader(500)
			return
		}

		c.metric.Histogram.With(prometheus.Labels{"code": "400", "method": "GET", "type": "findAll", "service": "category"}).Observe(time.Since(time.Now()).Seconds())
		w.WriteHeader(400)
		w.Write(response)
		return
	}

	response, err := json.Marshal(web.Success(data, "category found"))
	if err != nil {
		w.WriteHeader(500)
		return
	}

	c.metric.Histogram.With(prometheus.Labels{"code": "200", "method": "GET", "type": "findAll", "service": "category"}).Observe(time.Since(time.Now()).Seconds())
	w.WriteHeader(200)
	w.Write(response)
	return
}

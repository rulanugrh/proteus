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

type CommentInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
	FindUID(w http.ResponseWriter, r *http.Request)
	FindPID(w http.ResponseWriter, r *http.Request)
}

type comment struct {
	service service.CommentInterface
	metric *pkg.Metric
}

func CommentHandler(service service.CommentInterface, metric *pkg.Metric) CommentInterface {
	return &comment{service: service, metric: metric}
}

func(c *comment) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.Comment
	token := r.Header.Get("Authorization")
	claim, err := middleware.CheckToken(token)
	if err != nil {
		response, err := json.Marshal(web.Unauthorized(err.Error()))
		if err != nil {
			w.WriteHeader(500)
			return
		}

		c.metric.Histogram.With(prometheus.Labels{"code": "401", "method": "POST", "type": "create", "service": "comment"}).Observe(time.Since(time.Now()).Seconds())

		w.WriteHeader(401)
		w.Write(response)
		return
	}

	req.UserID = claim.ID
	req.Username = claim.Username
	req.Avatar = claim.Avatar
	req.RoleID = claim.RoleID
	
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	data, err := c.service.Create(req)
	if err != nil {
		response, err := json.Marshal(web.BadRequest(err.Error()))
		if err != nil {
			w.WriteHeader(500)
			return
		}

		c.metric.Histogram.With(prometheus.Labels{"code": "400", "method": "POST", "type": "create", "service": "comment"}).Observe(time.Since(time.Now()).Seconds())

		w.WriteHeader(400)
		w.Write(response)
		return
	}

	response, err := json.Marshal(web.Created(data, "success create domment"))
	if err != nil {
		w.WriteHeader(500)
		return
	}

	c.metric.Histogram.With(prometheus.Labels{"code": "201", "method": "POST", "type": "create", "service": "comment"}).Observe(time.Since(time.Now()).Seconds())
	c.metric.Counter.With(prometheus.Labels{"type": "create", "service": "category"}).Inc()
	w.WriteHeader(201)
	w.Write(response)
	return
}

func(c *comment) FindUID(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	claim, err := middleware.CheckToken(token)
	if err != nil {
		response, err := json.Marshal(web.Unauthorized(err.Error()))
		if err != nil {
			w.WriteHeader(500)
			return
		}

		c.metric.Histogram.With(prometheus.Labels{"code": "401", "method": "GET", "type": "findUID", "service": "comment"}).Observe(time.Since(time.Now()).Seconds())

		w.WriteHeader(401)
		w.Write(response)
		return
	}

	data, err := c.service.FindUID(claim.ID)
	if err != nil {
		response, err := json.Marshal(web.BadRequest(err.Error()))
		if err != nil {
			w.WriteHeader(500)
			return
		}

		c.metric.Histogram.With(prometheus.Labels{"code": "400", "method": "GET", "type": "findUID", "service": "comment"}).Observe(time.Since(time.Now()).Seconds())
		w.WriteHeader(400)
		w.Write(response)
		return
	}

	response, err := json.Marshal(web.Success(data, "comment in this user id found"))
	if err != nil {
		w.WriteHeader(500)
		return
	}

	c.metric.Histogram.With(prometheus.Labels{"code": "200", "method": "GET", "type": "findUID", "service": "comment"}).Observe(time.Since(time.Now()).Seconds())
	w.WriteHeader(200)
	w.Write(response)
	return
}

func(c *comment) FindPID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/comment/product/"))
	if err != nil {
		w.WriteHeader(500)
		return
	}

	data, err := c.service.FindUID(uint(id))
	if err != nil {
		response, err := json.Marshal(web.BadRequest(err.Error()))
		if err != nil {
			w.WriteHeader(500)
			return
		}
		c.metric.Histogram.With(prometheus.Labels{"code": "401", "method": "GET", "type": "findPID", "service": "comment"}).Observe(time.Since(time.Now()).Seconds())

		w.WriteHeader(400)
		w.Write(response)
		return
	}

	response, err := json.Marshal(web.Success(data, "comment in this product ID found"))
	if err != nil {
		w.WriteHeader(500)
		return
	}

	c.metric.Histogram.With(prometheus.Labels{"code": "200", "method": "GET", "type": "findPID", "service": "comment"}).Observe(time.Since(time.Now()).Seconds())

	w.WriteHeader(200)
	w.Write(response)
	return
}

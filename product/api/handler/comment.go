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
)

type CommentInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
	FindUID(w http.ResponseWriter, r *http.Request)
	FindPID(w http.ResponseWriter, r *http.Request)
}

type comment struct {
	service service.CommentInterface
}

func CommentHandler(service service.CommentInterface) CommentInterface {
	return &comment{service: service}
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

		w.WriteHeader(400)
		w.Write(response)
		return
	}

	response, err := json.Marshal(web.Created(data, "success create domment"))
	if err != nil {
		w.WriteHeader(500)
		return
	}

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

		w.WriteHeader(400)
		w.Write(response)
		return
	}

	response, err := json.Marshal(web.Success(data, "comment in this user id found"))
	if err != nil {
		w.WriteHeader(500)
		return
	}

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

		w.WriteHeader(400)
		w.Write(response)
		return
	}

	response, err := json.Marshal(web.Success(data, "comment in this product ID found"))
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	w.Write(response)
	return
}

package service

import (
	"strconv"

	"github.com/rulanugrh/tokoku/product/internal/entity/domain"
	"github.com/rulanugrh/tokoku/product/internal/entity/web"
	"github.com/rulanugrh/tokoku/product/internal/repository"
	"github.com/rulanugrh/tokoku/product/pkg"
)

type CommentInterface interface {
	Create(req domain.Comment) (*web.Comment, error)
	FindUID(id uint) (*[]web.Comment, error)
	FindPID(id uint) (*[]web.Comment, error)
}

type comment struct {
	repository repository.CommentInterface
	log pkg.ILogrus
}

func CommentService(repository repository.CommentInterface) CommentInterface {
	return &comment{repository: repository, log: pkg.Logrus()}
}

func (c *comment) Create(req domain.Comment) (*web.Comment, error) {
	data, err := c.repository.Create(req)
	if err != nil {
		c.log.Record("/api/comment/create", 500, "POST").Error(err.Error())
		return nil, err
	}

	response := web.Comment{
		Rate:     data.Rate,
		Product:  data.Product.Name,
		Comment:  data.Comment,
		Avatar:   data.Avatar,
		Username: data.Username,
	}

	c.log.Record("/api/comment/create", 200, "POST").Info("success creaet comment")
	return &response, nil
}

func (c *comment) FindUID(id uint) (*[]web.Comment, error) {
	data, err := c.repository.FindByUserID(id)
	if err != nil {
		c.log.Record("/api/comment/get", 500, "GET").Error(err.Error())
		return nil, err
	}

	var response []web.Comment
	for _, v := range *data {
		result := web.Comment{
			Product:  v.Product.Name,
			Comment:  v.Comment,
			Rate:     v.Rate,
			Username: v.Username,
			Avatar:   v.Avatar,
		}

		response = append(response, result)
	}

	c.log.Record("/api/comment/get", 200, "GET").Info("success get all comment with this userID "+strconv.Itoa(int(id)))
	return &response, nil
}

func (c *comment) FindPID(id uint) (*[]web.Comment, error) {
	data, err := c.repository.FindByProductID(id)
	if err != nil {
		c.log.Record("/api/comment/product/"+strconv.Itoa(int(id)), 500, "GET").Error(err.Error())
		return nil, err
	}

	var response []web.Comment
	for _, v := range *data {
		result := web.Comment{
			Product:  v.Product.Name,
			Comment:  v.Comment,
			Rate:     v.Rate,
			Username: v.Username,
			Avatar:   v.Avatar,
		}

		response = append(response, result)
	}

	c.log.Record("/api/comment/product/"+strconv.Itoa(int(id)), 200, "GET").Info("success get comment with this product ID")
	return &response, nil
}

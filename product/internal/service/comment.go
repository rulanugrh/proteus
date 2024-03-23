package service

import (
	"github.com/rulanugrh/tokoku/product/internal/entity/domain"
	"github.com/rulanugrh/tokoku/product/internal/entity/web"
	"github.com/rulanugrh/tokoku/product/internal/repository"
)

type CommentInterface interface {
	Create(req domain.Comment) (*web.Comment, error)
	FindUID(id uint) (*[]web.Comment, error)
	FindPID(id uint) (*[]web.Comment, error)
}

type comment struct {
	repository repository.CommentInterface
}

func CommentService(repository repository.CommentInterface) CommentInterface {
	return &comment{repository: repository}
}

func(c *comment) Create(req domain.Comment) (*web.Comment, error) {
	data, err := c.repository.Create(req)
	if err != nil {
		return nil, err
	}

	response := web.Comment{
		Rate: data.Rate,
		Product: data.Product.Name,
		Comment: data.Comment,
		Avatar: data.Avatar,
		Username: data.Username,
	}

	return &response, nil
}

func(c *comment) FindUID(id uint) (*[]web.Comment, error) {
	data, err := c.repository.FindByUserID(id)
	if err != nil {
		return nil, err
	}

	var response []web.Comment
	for _, v := range *data {
		result := web.Comment{
			Product: v.Product.Name,
			Comment: v.Comment,
			Rate: v.Rate,
			Username: v.Username,
			Avatar: v.Avatar,
		}

		response = append(response, result)
	}

	return &response, nil
}

func(c *comment) FindPID(id uint) (*[]web.Comment, error) {
	data, err := c.repository.FindByProductID(id)
	if err != nil {
		return nil, err
	}

	var response []web.Comment
	for _, v := range *data {
		result := web.Comment{
			Product: v.Product.Name,
			Comment: v.Comment,
			Rate: v.Rate,
			Username: v.Username,
			Avatar: v.Avatar,
		}

		response = append(response, result)
	}

	return &response, nil
}
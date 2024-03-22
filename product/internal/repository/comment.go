package repository

import (
	"github.com/rulanugrh/tokoku/product/internal/config"
	"github.com/rulanugrh/tokoku/product/internal/entity/domain"
)

type CommentInterface interface {
	Create(req domain.Comment) (*domain.Comment, error)
  FindByProductID(id uint) (*[]domain.Comment, error)
}

type comment struct {
  client config.Database
}

func CommentRepository(client config.Database) CommentInterface {
  return &comment{client: client}
}

func (c *comment) Create(req domain.Comment) (*domain.Comment, error) {
  err := c.client.DB.Create(&req).Error
  if err != nil {
    return nil, err
  }

  err = c.client.DB.Preload("Product").Find(&req).Error
  if err != nil {
    return nil, err
  }

  errAppend := c.client.DB.Model(&req.Product).Association("Product").Append(&req)
  if errAppend != nil {
    return nil, errAppend
  }

  return &req, nil
}

func (c *comment) FindByProductID(id uint) (*[]domain.Comment, error) {
  var model []domain.Comment
  err := c.client.DB.Preload("Product").Where("product_id = ?", id).Find(&model).Error
  if err != nil {
    return nil, err
  }

  return &model, nil
}

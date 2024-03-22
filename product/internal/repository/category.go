package repository

import (
	"github.com/rulanugrh/tokoku/product/internal/config"
	"github.com/rulanugrh/tokoku/product/internal/entity/domain"
)

type CategoryInterface interface {
	Create(req domain.Category) (*domain.Category, error)
  FindID(id uint) (*domain.Category, error)
  FindAll() (*[]domain.Category, error)
}

type category struct {
  client config.Database
}

func CategoryRepository(client config.Database) CategoryInterface {
  return &category{client: client}
}

func (c *category) Create(req domain.Category) (*domain.Category, error) {
  err := c.client.DB.Create(&req).Error
  if err != nil {
    return nil, err
  }

  err = c.client.DB.Preload("Product").Find(&req).Error
  if err != nil {
    return nil, err
  }

  return &req, nil
}

func (c *category) FindID(id uint) (*domain.Category, error) {
  var model domain.Category

  err := c.client.DB.Preload("Product").Where("id = ?", id).Find(&model).Error
  if err != nil {
    return nil, err
  }

  return &model, nil
}

func (c *category) FindAll() (*[]domain.Category, error) {
  var model []domain.Category
  err := c.client.DB.Preload("Product").Find(&model).Error

  if err != nil {
    return nil, err
  }

  return &model, nil
}

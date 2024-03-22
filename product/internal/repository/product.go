package repository

import (
	"github.com/rulanugrh/tokoku/product/internal/config"
	"github.com/rulanugrh/tokoku/product/internal/entity/domain"
)

type ProductInterface interface {
	Create(req domain.Product) (*domain.Product, error)
  FindID(id uint) (*domain.Product, error)
  FindAll() (*[]domain.Product, error)
  Update(id uint, req domain.Product) (*domain.Product, error)
}

type product struct {
  client config.Database
}

func ProductRepository(client config.Database) ProductInterface {
  return &product{client: client}
}

func (p *product) Create(req domain.Product) (*domain.Product, error) {
  err := p.client.DB.Create(&req).Error

  if err != nil {
    return nil, err
  }

  err = p.client.DB.Preload("Category").Preload("Comment").Find(&req).Error
  if err != nil {
    return nil, err
  }

  errAppend := p.client.DB.Model(&req.Category).Association("Product").Append(&req)
  if errAppend != nil {
    return nil, errAppend
  }

  return &req, nil
}

func (p *product) FindID(id uint) (*domain.Product, error) {
  var model domain.Product

  err := p.client.DB.Preload("Comment").Preload("Category").Where("id = ?", id).Find(&model).Error

  if err != nil {
    return nil, err
  }

  return &model, nil

}

func (p *product) FindAll() (*[]domain.Product, error) {
  var model []domain.Product

  err := p.client.DB.Preload("Comment").Preload("Category").Find(&model).Error
  if err != nil {
    return nil, err
  }

  return &model, nil
}

func (p *product) Update(id uint, req domain.Product) (*domain.Product, error) {
  var model domain.Product

  err := p.client.DB.Model(&req).Where("id = ?", id).Updates(&model).Error
  if err != nil {
    return nil, err
  }

  preload := p.client.DB.Preload("Comment").Preload("Category").Where("id = ?", id).Find(&model).Error

  if preload != nil {
    return nil, err
  }

  return &model, nil
}

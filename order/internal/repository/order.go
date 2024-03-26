package repository

import (
	"github.com/rulanugrh/order/internal/config"
	"github.com/rulanugrh/order/internal/entity"
	"github.com/rulanugrh/order/internal/util"
)

type OrderInterface interface {
	Create(req entity.Order) (*entity.Order, error)
	Checkout(uuid string)  (*entity.Order, error)
	Update(uuid string, status string) error
}

type order struct {
	client *config.Postgres
}

func OrderRepository(client *config.Postgres) OrderInterface {
	return &order{client: client}
}

func(o *order) Create(req entity.Order) (*entity.Order, error) {
	req.Status = "not paid"
	req.UUID = util.GenerateUUID()
	err := o.client.DB.Create(&req).Error
	if err != nil {
		return nil, err
	}

	preload := o.client.DB.Preload("Product").Find(&req).Error
	if preload != nil {
		return nil, err
	}

	return &req, nil
}

func(o *order) Checkout(uuid string)  (*entity.Order, error) {
	var model entity.Order
	err := o.client.DB.Preload("Product").Where("uuid = ?", uuid).Find(&model).Error
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (o *order) Update(uuid string, status string) error {
	err := o.client.DB.Model(&entity.Order{}).Where("uuid = ?", uuid).Update("status", status).Error
	if err != nil {
		return err
	}

	return nil
}
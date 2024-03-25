package repository

import (
	"github.com/rulanugrh/order/internal/config"
	"github.com/rulanugrh/order/internal/entity"
)

type OrderInterface interface {
	Create(req entity.Order) (*entity.Order, error)
	Pay(uuid string, status string) error
}

type order struct {
	client *config.Postgres
}

func OrderRepository(client *config.Postgres) OrderInterface {
	return &order{client: client}
}

func(o *order) Create(req entity.Order) (*entity.Order, error) {
	req.Status = "not paid"
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

func(o *order) Pay(uuid string, status string)  error {
	err := o.client.DB.Model(&entity.Order{}).Where("uuid = ?", uuid).Update("status", status).Error
	if err != nil {
		return err
	}

	return nil
}
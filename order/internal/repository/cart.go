package repository

import (
	"github.com/rulanugrh/order/internal/config"
	"github.com/rulanugrh/order/internal/entity"
	"github.com/rulanugrh/order/internal/util"
)

type CartInterface interface {
	AddToCart(req entity.Cart) error
	ListCart(userID uint) (*[]entity.Cart, error)
	ProcessCart(id uint, updates entity.Updates) (*entity.Order, error)
	Update(id uint, req entity.Cart) error
	Delete(id uint) error
}

type cart struct {
	client *config.Postgres
}

func CartRepository(client *config.Postgres) CartInterface {
	return &cart{client: client}
}

func(c *cart) AddToCart(req entity.Cart) error {
	err := c.client.DB.Create(&req).Error
	if err != nil {
		return err
	}

	return nil
}

func(c *cart) ListCart(userID uint) (*[]entity.Cart, error) {
	var response []entity.Cart
	find := c.client.DB.Where("user_id = ?", userID).Preload("Product").Find(&response)
	if find.RowsAffected == 0 {
		return nil, find.Error
	}

	return &response, nil
}

func(c *cart) ProcessCart(id uint, updates entity.Updates) (*entity.Order, error)  {
	var model entity.Cart
	find := c.client.DB.Where("id = ?", id).Find(&model)
	
	if find.RowsAffected == 0 {
		return nil, find.Error
	}

	var order entity.Order
	order.UUID = util.GenerateUUID()
	order.UserID = model.UserID
	order.ProductID = model.ProductID
	order.Quantity = model.Quantity
	order.Address = updates.Address
	order.MethodPayment = updates.MethodType

	err_create := c.client.DB.Create(&order).Error
	if err_create != nil {
		return nil, err_create
	}

	err_preload := c.client.DB.Preload("Product").Find(&order).Error
	if err_preload != nil {
		return nil, err_preload
	}

	return &order, nil
}

func(c *cart) Update(id uint, req entity.Cart) error {
	var updates entity.Cart

	err := c.client.DB.Model(&req).Where("id = ?", id).Updates(&updates).Error
	if err != nil {
		return err
	}

	return nil
}

func(c *cart) Delete(id uint) error {
	var model entity.Cart

	err := c.client.DB.Where("id = ?", id).Delete(&model).Error
	if err != nil {
		return err
	}

	return nil
}

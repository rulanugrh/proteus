package repository

import (
	"github.com/rulanugrh/order/internal/config"
	"github.com/rulanugrh/order/internal/entity"
	"github.com/rulanugrh/order/internal/util"
	"github.com/rulanugrh/order/internal/util/constant"
	"github.com/rulanugrh/order/pkg/logger"
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
	log logger.ILogrus
}

func CartRepository(client *config.Postgres) CartInterface {
	return &cart{client: client, log: logger.Logrus()}
}

func (c *cart) AddToCart(req entity.Cart) error {
	err := c.client.DB.Create(&req).Error
	if err != nil {
		c.log.RecordDB("AddToCart", err.Error()).Error("error create cart")
		return constant.InternalServerError("error while created cart", err)
	}

	c.log.RecordDB("AddToCart", nil).Info("success create cart")
	return nil
}

func (c *cart) ListCart(userID uint) (*[]entity.Cart, error) {
	var response []entity.Cart
	find := c.client.DB.Where("user_id = ?", userID).Preload("Product").Find(&response)
	if find.RowsAffected == 0 {
		c.log.RecordDB("ListCart", find.Error.Error()).Error("user id not found")
		return nil, constant.NotFound("sorry your list cart with this id not found")
	}

	c.log.RecordDB("ListCart", nil).Info("cart found")
	return &response, nil
}

func (c *cart) ProcessCart(id uint, updates entity.Updates) (*entity.Order, error) {
	var model entity.Cart
	find := c.client.DB.Where("id = ?", id).Find(&model)

	if find.RowsAffected == 0 {
		c.log.RecordDB("ProcessCart", find.Error.Error()).Error("cart not found")
		return nil, constant.NotFound("sorry cart with this id not found")
	}

	var order entity.Order
	order.UUID = util.GenerateUUID()
	order.UserID = model.UserID
	order.ProductID = model.ProductID
	order.Quantity = model.Quantity
	order.Address = updates.Address
	order.MethodPayment = updates.MethodType
	order.RequestCurreny = updates.RequestCurreny
	order.ChannelCode = updates.ChannelCode
	order.MobilePhone = updates.MobilePhone

	err_create := c.client.DB.Create(&order).Error
	if err_create != nil {
		c.log.RecordDB("ProcessCart", err_create.Error()).Error("errror while create order")
		return nil, constant.InternalServerError("sorry cannot create order", err_create)
	}

	c.log.RecordDB("ProcessCart", nil).Info("success proccess cart")
	return &order, nil
}

func (c *cart) Update(id uint, req entity.Cart) error {
	var updates entity.Cart

	err := c.client.DB.Model(&req).Where("id = ?", id).Updates(&updates).Error
	if err != nil {
		c.log.RecordDB("Update", err.Error()).Error("error while update cart")
		return constant.InternalServerError("error while update cart", err)
	}

	c.log.RecordDB("Update", nil).Info("success update cart")
	return nil
}

func (c *cart) Delete(id uint) error {
	var model entity.Cart

	err := c.client.DB.Where("id = ?", id).Delete(&model).Error
	if err != nil {
		c.log.RecordDB("Delete", err.Error()).Error("error delete cart")
		return constant.InternalServerError("error cannot deleted data", err)
	}

	c.log.RecordDB("Delete", nil).Info("success delete cart")
	return nil
}

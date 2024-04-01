package repository

import (
	"github.com/rulanugrh/order/internal/config"
	"github.com/rulanugrh/order/internal/entity"
	"github.com/rulanugrh/order/internal/util"
	"github.com/rulanugrh/order/internal/util/constant"
	"github.com/rulanugrh/order/pkg/logger"
)

type OrderInterface interface {
	Create(req entity.Order) (*entity.Order, error)
	Checkout(uuid string) (*entity.Order, error)
	Update(uuid string, status string) error
	SaveTransaction(req entity.Transaction) error
}

type order struct {
	client *config.Postgres
	log    logger.ILogrus
}

func OrderRepository(client *config.Postgres) OrderInterface {
	return &order{client: client, log: logger.Logrus()}
}

func (o *order) Create(req entity.Order) (*entity.Order, error) {
	req.Status = "not paid"
	req.UUID = util.GenerateUUID()

	err := o.client.DB.Create(&req).Error
	if err != nil {
		o.log.RecordDB("Create", err.Error()).Error("error while create order")
		return nil, constant.InternalServerError("sorry cannot create order", err)
	}

	o.log.RecordDB("Create", nil).Info("success create order")
	return &req, nil
}

func (o *order) Checkout(uuid string) (*entity.Order, error) {
	var model entity.Order
	err := o.client.DB.Where("uuid = ?", uuid).Find(&model).Error
	if err != nil {
		o.log.RecordDB("Checkout", err.Error()).Error("error find order")
		return nil, constant.InternalServerError("something error while checkout", err)
	}

	o.log.RecordDB("Checkout", nil).Info("data found")
	return &model, nil
}

func (o *order) Update(uuid string, status string) error {
	err := o.client.DB.Model(&entity.Order{}).Where("uuid = ?", uuid).Update("status", status).Error
	if err != nil {
		o.log.RecordDB("Update", err.Error()).Error("error while updating order")
		return constant.InternalServerError("something error while update data", err)
	}

	o.log.RecordDB("Update", nil).Info("success update data")
	return nil
}

func (o *order) SaveTransaction(req entity.Transaction) error {
	err := o.client.DB.Create(&req).Error
	if err != nil {
		o.log.RecordDB("SaveTransaction", err.Error()).Error("error while create transaction")
		return constant.InternalServerError("error while create transaction", err)
	}

	o.log.RecordDB("SaveTransaction", nil).Info("success save into db")
	return nil
}

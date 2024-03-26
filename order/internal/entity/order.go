package entity

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UUID          string  `json:"uuid" form:"uuid"`
	UserID        uint    `json:"user_id" form:"user_id"`
	ProductID     uint    `json:"product_id" form:"product_id" validate:"required"`
	Quantity      uint    `json:"quantity" form:"quantity" validate:"required"`
	Status        string  `json:"status" form:"status"`
	MethodPayment string  `json:"method_payment" form:"method_payment validate:required"`
	Address       string  `json:"address" form:"address" validate:"required"`
	Product       Product `json:"product" form:"product" gorm:"foreignKey:ProductID;references:ID"`
}

package entity

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UUID      string  `json:"uuid" form:"uuid"`
	UserID    uint    `json:"user_id" form:"user_id"`
	ProductID uint    `json:"product_id" form:"product_id" validate:"required"`
	Product   Product `json:"product" form:"product" gorm:"foreignKey:ProductID;references:ID"`
	Quantity  uint    `json:"quantity" form:"quantity" validate:"required"`
	Paid      bool    `json:"paid" form:"paid"`
}

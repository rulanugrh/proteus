package entity

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID    uint    `json:"user_id" form:"user_id"`
	ProductID uint    `json:"product_id" form:"product_id" validate:"required"`
	Quantity  uint    `json:"quantity" form:"quantity" validate:"required"`
}

type Updates struct {
	MethodType     string `json:"method_type" form:"method_type"`
	Address        string `json:"address" form:"address"`
	RequestCurreny string `json:"request_currency" form:"request_currenty"`
	ChannelCode    string  `json:"channel_code" form:"channel_code" validate:"required"`
	MobilePhone    string  `json:"mobil_phone" form:"mobile_phone"`
}

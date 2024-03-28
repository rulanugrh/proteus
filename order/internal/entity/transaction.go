package entity

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	OrderID                uint    `json:"order_id" form:"order_id"`
	OrderUUID              string  `json:"order_uuid" form:"order_uuid"`
	MethodPayment          string  `json:"method_payment"`
	Status                 string  `json:"status"`
	PaymentCreated         string  `json:"payment_created"`
	PaymentUpdated         string  `json:"payment_updated"`
	PaymentRequestCurrency string  `json:"payment_currency"`
	PaymentID              string  `json:"payment_id"`
	Amount                 float64 `json:"amount"`
	ReferenceID            string  `json:"reference_id"`
}

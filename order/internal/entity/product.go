package entity

import "time"

type Product struct {
	ID           uint      `json:"id" form:"id"`
	Name         string    `json:"name" form:"name"`
	Description  string    `json:"desc" form:"desc"`
	Price        uint32    `json:"price" form:"price"`
	QtyAvailable uint64    `json:"qty_available" form:"qty_available"`
	CreateAt     time.Time `json:"create_at" form:"create_at"`
	UpdateAt     time.Time `json:"update_at" form:"update_at"`
}

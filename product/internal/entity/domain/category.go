package domain

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name        string    `json:"name" form:"name"`
	Description string    `json:"desc" form:"desc"`
	Product     []Product `json:"product" form:"product"`
}

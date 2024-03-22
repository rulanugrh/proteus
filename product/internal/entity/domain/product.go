package domain

import "gorm.io/gorm"

type Product struct {
	gorm.Model

  Name string `json:"name" form:"name" validate:"required"`
  Description string `json:"desc" form:"desc" validate:"required"`
  Price uint32 `json:"price" form:"price" validate:"required"`
  QtyAvailable uint64 `json:"qty_available" form:"qty_available"`
  QtyOn uint64 `json:"qty_on" form:"qty_on"`
  QtyReserved uint64 `json:"qty_reserved" form:"qty_reserved" validate:"requird"`
  CategoryID uint `json:"category_id" form:"category_id"`
  Category Category `json:"category" form:"category" gorm:"foreignKey:CategoryID;reference:ID"`
  Comment []Comment `json:"comment" form:"comment" gorm:"many2many:review"`
}

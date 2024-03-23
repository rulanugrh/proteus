package domain

type Comment struct {
	UserID    uint    `json:"user_id" form:"user_id"`
	ProductID uint    `json:"product_id" form:"product_id"`
	Comment   string  `json:"comment" form:"comment"`
	Rate      int8    `json:"rate" form:"rate"`
	Product   Product `json:"product" form:"product" gorm:"foreignKey:ProductID;references:ID"`
}

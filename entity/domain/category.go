package domain

type Category struct {
	ID       uint
	Nama     string
	Products []Product `gorm:"foreignKey:CategoryName;references:Nama"`
}

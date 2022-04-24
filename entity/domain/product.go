package domain

type Product struct {
	ID           uint
	Nama         string
	Price        int
	CategoryName string
}

type PaginationProduct struct {
	Limit int
	Page  int
	Sort  string
}

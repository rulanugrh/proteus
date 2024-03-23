package web

type Product struct {
	ID          uint   `json:"id"`
	Name        string `json:"name" form:"name"`
	Description string `json:"desc" form:"desc"`
	Price       uint32 `json:"price" form:"price"`
	Available   uint64 `json:"available" form:"available"`
	Reserved    uint64 `json:"reserved" form:"reserved"`
	Category    string `json:"category" form:"category"`
}

type GetProduct struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name" form:"name"`
	Description string    `json:"desc" form:"desc"`
	Price       uint32    `json:"price" form:"price"`
	Available   uint64    `json:"available" form:"available"`
	Reserved    uint64    `json:"reserved" form:"reserved"`
	Category    string    `json:"category" form:"category"`
	Comment     []Comment `json:"comment" form:"comment"`
}

type Category struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name" form:"name"`
	Description string    `json:"description"`
	Product     []Product `json:"product" form:"product"`
}

type Comment struct {
	Rate    int8   `json:"rate"`
	Comment string `json:"comment"`
	Product string `json:"product"`
}

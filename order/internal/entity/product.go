package entity

type Product struct {
	ID          uint   `json:"id"`
	Name        string `json:"name" form:"name"`
	Description string `json:"desc" form:"desc"`
	Price       uint32 `json:"price" form:"price"`
	Available   uint64 `json:"available" form:"available"`
	Reserved    uint64 `json:"reserved" form:"reserved"`
	Category    string `json:"category" form:"category"`
}

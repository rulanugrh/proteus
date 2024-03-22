package web

type Product struct {
  ID uint `json:"id"`
  Name string `json:"name" form:"name"`
  Description string `json:"desc" form:"desc"`
  Price uint32 `json:"price" form:"price"`
  Available uint32 `json:"available" form:"available"`
  Reserved uint32 `json:"reserved" form:"reserved"`
  Category string `json:"category" form:"category"`
  Comment []Comment `json:"comment" form:"comment"`
}

type GetProduct struct {
  Name string `json:"name" form:"name"`
  Description string `json:"desc" form:"desc"`
  Price uint32 `json:"price" form:"price"`
  Available uint32 `json:"available" form:"available"`
  Reserved uint32 `json:"reserved" form:"reserved"`
  Category string `json:"category" form:"category"`
}

type Category struct {
  ID uint `json:"id"`
  Name string `json:"name" form:"name"`
  Description string `json:"description"`
  Product []GetProduct `json:"product" form:"product"`
}

type Comment struct {
  Username string `json:"username"`
  Avatar string `json:"avatar"`
  Comment string `json:"comment"`
  Product string `json:"product"`
}

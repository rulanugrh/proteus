package app

import (
	productRepo "github.com/ItsArul/TokoKu/controller/interfaces"
	"github.com/gin-gonic/gin"
)

func Run(product productRepo.ProductController) {
	server := gin.Default()

	route := server.Group("/api/product")
	{
		route.POST("/", product.Create())
		route.GET("/:id", product.FindById())
		route.GET("/", product.FindAll())
		route.PUT("/:id", product.Update())
		route.DELETE("/:id", product.Delete())
	}

	server.Run(":8080")
}

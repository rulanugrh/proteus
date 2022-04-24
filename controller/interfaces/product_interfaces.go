package interfaces

import "github.com/gin-gonic/gin"

type ProductController interface {
	Create() gin.HandlerFunc
	Update() gin.HandlerFunc
	FindById() gin.HandlerFunc
	FindAll() gin.HandlerFunc
	Delete() gin.HandlerFunc
}

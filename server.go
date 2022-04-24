package main

import (
	"github.com/ItsArul/TokoKu/app"
	"github.com/ItsArul/TokoKu/controller"
	"github.com/ItsArul/TokoKu/entity/domain"
	"github.com/ItsArul/TokoKu/repository"
	"github.com/ItsArul/TokoKu/services"
	"github.com/ItsArul/TokoKu/utilities"
)

func main() {
	DB := utilities.GetConnection()
	DB.AutoMigrate(&domain.Product{}, &domain.Category{})

	productRepo := repository.StartProductRepository()
	productService := services.StartProductServices(productRepo)
	productControll := controller.StartProductController(productService)

	app.Run(productControll)
}

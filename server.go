package main

import (
	"github.com/ItsArul/TokoKu/app"
	"github.com/ItsArul/TokoKu/config"
	"github.com/ItsArul/TokoKu/controller"
	"github.com/ItsArul/TokoKu/entity/domain"
	"github.com/ItsArul/TokoKu/repository"
	"github.com/ItsArul/TokoKu/services"
)

func main() {
	DB := config.GetConnect()
	DB.AutoMigrate(&domain.Product{}, &domain.Category{})

	productRepo := repository.StartProductRepository()
	productService := services.StartProductServices(productRepo)
	productControll := controller.StartProductController(productService)

	app.Run(productControll)
}

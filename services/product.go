package services

import (
	"context"

	"github.com/ItsArul/TokoKu/entity/domain"
	"github.com/ItsArul/TokoKu/repository/interfaces"
	productservices "github.com/ItsArul/TokoKu/services/interfaces"
	"github.com/ItsArul/TokoKu/utilities"
)

type productinterfaces struct {
	productinter interfaces.ProductRepo
}

func StartProductServices(product interfaces.ProductRepo) productservices.ProducInterfaces {
	return &productinterfaces{productinter: product}
}

func (prod *productinterfaces) Create(ctx context.Context, product domain.Product) (domain.Product, error) {
	productdomain, err := prod.productinter.Create(product, ctx)
	if err != nil {
		utilities.Error(err)
	}
	return productdomain, nil
}

func (prod *productinterfaces) FindById(ctx context.Context, id uint) (domain.Product, error) {
	productdomain, err := prod.productinter.FindById(ctx, id)
	if err != nil {
		utilities.Error(err)
	}

	return productdomain, nil
}

func (prod *productinterfaces) Update(ctx context.Context, id uint, product domain.Product) (domain.Product, error) {
	productdomain, err := prod.productinter.Update(ctx, id, product)
	if err != nil {
		utilities.Error(err)
	}

	return productdomain, nil
}

func (prod *productinterfaces) Delete(ctx context.Context, id uint) error {
	err := prod.productinter.Delete(ctx, id)
	if err != nil {
		utilities.Error(err)
	}

	return nil
}

func (prod *productinterfaces) FindAll(ctx context.Context) ([]domain.Product, error) {
	productDomain := prod.productinter.FindAll(ctx)

	return productDomain, nil

}

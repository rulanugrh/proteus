package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/ItsArul/TokoKu/entity/domain"
	"github.com/ItsArul/TokoKu/repository/interfaces"
	"github.com/ItsArul/TokoKu/utilities"
)

type productrepository struct {
}

func StartProductRepository() interfaces.ProductRepo {
	return &productrepository{}
}

func (p *productrepository) Create(product domain.Product, ctx context.Context) (domain.Product, error) {
	tx := utilities.Dbconn.WithContext(ctx)
	err := tx.Create(&product).Error
	if err != nil {
		return product, errors.New("cannot create product")
	}

	return product, nil
}

func (p *productrepository) FindById(ctx context.Context, id uint) (domain.Product, error) {
	var product domain.Product

	tx := utilities.Dbconn.WithContext(ctx)
	err := tx.First(&product, id).Error
	if err != nil {
		return product, errors.New("cannot find product by id")
	}

	return product, nil

}

func (p *productrepository) FindAll(ctx context.Context, pagination domain.PaginationProduct) []domain.Product {
	var product []domain.Product

	tx := utilities.Dbconn.WithContext(ctx)
	offset := (pagination.Page - 1) * pagination.Limit
	build := tx.Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	err := build.Find(&product).Error
	if err != nil {
		fmt.Println("cannot find all product")
		panic(err)
	}

	return product
}

func (p *productrepository) Update(ctx context.Context, id uint, product domain.Product) (domain.Product, error) {
	var products domain.Product

	tx := utilities.Dbconn.WithContext(ctx)
	err := tx.Model(&products).Where("id = ?", id).Updates(&product).Error
	if err != nil {
		return product, errors.New("cannoot update product")
	}

	return product, nil
}

func (p *productrepository) Delete(ctx context.Context, id uint) error {
	var product domain.Product

	tx := utilities.Dbconn.WithContext(ctx)
	err := tx.Delete(&product).Where("id = ?", id).Error
	if err != nil {
		return errors.New("cannot delete product")
	}

	return nil
}

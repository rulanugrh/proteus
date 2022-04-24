package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/ItsArul/TokoKu/config"
	"github.com/ItsArul/TokoKu/entity/domain"
	"github.com/ItsArul/TokoKu/repository/interfaces"
)

type productrepository struct {
}

func StartProductRepository() interfaces.ProductRepo {
	return &productrepository{}
}

func (p *productrepository) Create(product domain.Product, ctx context.Context) (domain.Product, error) {
	tx := config.Dbconn.WithContext(ctx)
	err := tx.Create(&product).Error
	if err != nil {
		return product, errors.New("cannot create product")
	}

	return product, nil
}

func (p *productrepository) FindById(ctx context.Context, id uint) (domain.Product, error) {
	var product domain.Product

	tx := config.Dbconn.WithContext(ctx)
	err := tx.First(&product, id).Error
	if err != nil {
		return product, errors.New("cannot find product by id")
	}

	return product, nil

}

func (p *productrepository) FindAll(ctx context.Context) []domain.Product {
	var product []domain.Product

	tx := config.Dbconn.WithContext(ctx)
	err := tx.Find(&product).Error
	if err != nil {
		fmt.Println("cannot find all product")
		panic(err)
	}

	return product
}

func (p *productrepository) Update(ctx context.Context, id uint) (domain.Product, error) {
	var product domain.Product

	tx := config.Dbconn.WithContext(ctx)
	err := tx.Model(&product).Where("id = ?", id).Updates(&product).Error
	if err != nil {
		return product, errors.New("cannoot update product")
	}

	return product, nil
}

func (p *productrepository) Delete(ctx context.Context, id uint) error {
	var product domain.Product

	tx := config.Dbconn.WithContext(ctx)
	err := tx.Delete(&product).Where("id = ?", id).Error
	if err != nil {
		return errors.New("cannot delete product")
	}

	return nil
}

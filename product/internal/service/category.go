package service

import (
	"github.com/rulanugrh/tokoku/product/internal/entity/domain"
	"github.com/rulanugrh/tokoku/product/internal/entity/web"
	"github.com/rulanugrh/tokoku/product/internal/repository"
)

type CategoryInterface interface {
	Create(req domain.Category) (*web.Category, error)
	FindID(id uint) (*web.GetCategory, error)
	FindAll() (*[]web.GetCategory, error)
}

type category struct {
	repository repository.CategoryInterface
}

func CategoryService(repository repository.CategoryInterface) CategoryInterface {
	return &category{repository: repository}
}

func (c *category) Create(req domain.Category) (*web.Category, error) {
	data, err := c.repository.Create(req)
	if err != nil {
		return nil, err
	}

	response := web.Category{
		ID:          data.ID,
		Name:        data.Name,
		Description: data.Description,
	}

	return &response, nil
}

func (c *category) FindID(id uint) (*web.GetCategory, error) {
	data, err := c.repository.FindID(id)
	if err != nil {
		return nil, err
	}

	var products []web.Product
	for _, v := range data.Product {
		product := web.Product{
			ID:          v.ID,
			Available:   v.QtyAvailable,
			Reserved:    v.QtyReserved,
			Price:       v.Price,
			Name:        v.Name,
			Description: v.Description,
			Category:    v.Category.Name,
		}

		products = append(products, product)
	}

	response := web.GetCategory{
		ID:          data.ID,
		Name:        data.Name,
		Description: data.Description,
		Product:     products,
	}

	return &response, nil
}

func (c *category) FindAll() (*[]web.GetCategory, error) {
	data, err := c.repository.FindAll()
	if err != nil {
		return nil, err
	}

	var response []web.GetCategory
	var products []web.Product
	for _, v := range *data {
		for _, p := range v.Product {
			product := web.Product{
				ID:          p.ID,
				Price:       p.Price,
				Name:        p.Name,
				Category:    p.Category.Name,
				Available:   p.QtyAvailable,
				Reserved:    p.QtyReserved,
				Description: p.Description,
			}
			products = append(products, product)
		}

		result := web.GetCategory{
			ID:          v.ID,
			Description: v.Description,
			Name:        v.Name,
			Product:     products,
		}

		response = append(response, result)
	}

	return &response, nil
}

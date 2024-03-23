package service

import (
	"github.com/rulanugrh/tokoku/product/internal/entity/domain"
	"github.com/rulanugrh/tokoku/product/internal/entity/web"
	"github.com/rulanugrh/tokoku/product/internal/repository"
)

type ProductInterface interface {
	Create(req domain.Product) (*web.Product, error)
	FindID(id uint) (*web.GetProduct, error)
	FindAll() (*[]web.GetProduct, error)
	Update(id uint, req domain.Product) (*web.Product, error)
}

type product struct {
	repository repository.ProductInterface
}

func ProductService(repository repository.ProductInterface) ProductInterface {
	return &product{repository: repository}
}

func (p *product) Create(req domain.Product) (*web.Product, error) {
	data, err := p.repository.Create(req)
	if err != nil {
		return nil, err
	}

	response := web.Product{
		ID:          data.ID,
		Name:        data.Name,
		Description: data.Description,
		Price:       data.Price,
		Available:   data.QtyAvailable,
		Reserved:    data.QtyReserved,
		Category:    data.Category.Name,
	}

	return &response, nil
}

func (p *product) FindID(id uint) (*web.GetProduct, error) {
	data, err := p.repository.FindID(id)
	if err != nil {
		return nil, err
	}

	var comment []web.Comment
	for _, v := range data.Comment {
		result := web.Comment{
			Comment: v.Comment,
			Product: data.Name,
			Rate:    v.Rate,
		}

		comment = append(comment, result)
	}

	response := web.GetProduct{
		ID:          data.ID,
		Name:        data.Name,
		Description: data.Description,
		Price:       data.Price,
		Available:   data.QtyAvailable,
		Reserved:    data.QtyReserved,
		Category:    data.Category.Name,
		Comment:     comment,
	}

	return &response, nil
}

func (p *product) FindAll() (*[]web.GetProduct, error) {
	data, err := p.repository.FindAll()
	if err != nil {
		return nil, err
	}

	var response []web.GetProduct
	var comments []web.Comment
	for _, result := range *data {
		for _, c := range result.Comment {
			comment := web.Comment{
				Rate:    c.Rate,
				Product: c.Product.Name,
				Comment: c.Comment,
			}

			comments = append(comments, comment)
		}

		res := web.GetProduct{
			ID:          result.ID,
			Category:    result.Category.Name,
			Name:        result.Name,
			Description: result.Description,
			Price:       result.Price,
			Available:   result.QtyAvailable,
			Reserved:    result.QtyReserved,
			Comment:     comments,
		}

		response = append(response, res)
	}

	return &response, nil
}

func (p *product) Update(id uint, req domain.Product) (*web.Product, error) {
	data, err := p.repository.Update(id, req)
	if err != nil {
		return nil, err
	}

	response := web.Product{
		ID:        data.ID,
		Name:      data.Name,
		Price:     data.Price,
		Available: data.QtyAvailable,
		Reserved:  data.QtyReserved,
		Category:  data.Category.Name,
	}

	return &response, nil
}

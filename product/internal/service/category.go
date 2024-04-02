package service

import (
	"strconv"

	"github.com/rulanugrh/tokoku/product/internal/entity/domain"
	"github.com/rulanugrh/tokoku/product/internal/entity/web"
	"github.com/rulanugrh/tokoku/product/internal/repository"
	"github.com/rulanugrh/tokoku/product/pkg"
)

type CategoryInterface interface {
	Create(req domain.Category) (*web.Category, error)
	FindID(id uint) (*web.GetCategory, error)
	FindAll() (*[]web.GetCategory, error)
}

type category struct {
	repository repository.CategoryInterface
	log pkg.ILogrus
}

func CategoryService(repository repository.CategoryInterface) CategoryInterface {
	return &category{repository: repository, log: pkg.Logrus()}
}

func (c *category) Create(req domain.Category) (*web.Category, error) {
	data, err := c.repository.Create(req)
	if err != nil {
		c.log.Record("/api/category/create", 500, "POST").Error(err.Error())
		return nil, err
	}

	response := web.Category{
		ID:          data.ID,
		Name:        data.Name,
		Description: data.Description,
	}

	c.log.Record("/api/category/create", 200, "POST").Info("success create category")
	return &response, nil
}

func (c *category) FindID(id uint) (*web.GetCategory, error) {
	data, err := c.repository.FindID(id)
	if err != nil {
		c.log.Record("/api/category/find/"+strconv.Itoa(int(id)), 500, "GET").Error(err.Error())
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

	c.log.Record("/api/category/find/"+strconv.Itoa(int(id)), 200, "GET").Info("success get category by this id "+strconv.Itoa(int(id)))
	return &response, nil
}

func (c *category) FindAll() (*[]web.GetCategory, error) {
	data, err := c.repository.FindAll()
	if err != nil {
		c.log.Record("/api/category/get", 500, "GET").Error(err.Error())
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

	c.log.Record("/api/category/get", 200, "GET").Info("success get all category")
	return &response, nil
}

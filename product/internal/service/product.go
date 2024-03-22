package service

import (
	"github.com/rulanugrh/tokoku/product/internal/entity/domain"
	"github.com/rulanugrh/tokoku/product/internal/entity/web"
	"github.com/rulanugrh/tokoku/product/internal/repository"
)

type ProductInterface interface {
	Create(req domain.Product) (*web.Product, error)
  FindID(id uint) (*web.Product, error)
  FindAll() (*[]web.Product, error)
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
    ID: data.ID,
    Name: data.Name,
    Description: data.Description,
    Price: data.Price,
    Available: uint32(data.QtyAvailable),
    Reserved: uint32(data.QtyReserved),
    Category: data.Category.Name,
  }

  return &response, nil
}

func (p *product) FindID(id uint) (*web.Product, error) {
  data, err := p.repository.FindID(id)
  if err != nil {
    return nil, err
  }

  var comment []web.Comment
  for _, v := range data.Comment {
    result := web.Comment{
      Comment: v.Comment,
      Username: v.Product.Name,
      Product: data.Name,
      Avatar: data.Description,
    }

    comment = append(comment, result)
  }

  response := web.Product{
    ID: data.ID,
    Name: data.Name,
    Description: data.Description,
    Price: data.Price,
    Available: uint32(data.QtyAvailable),
    Reserved: uint32(data.QtyReserved),
    Category: data.Category.Name,
    Comment: comment,
  }

  return &response, nil
}


func (p *product) FindAll() (*[]web.Product, error) {
  data, err := p.repository.FindAll()
  if err != nil {
    return nil, err
  }

  var response []web.Product
  var comments []web.Comment
  for _, result := range *data {
    for _, c := range result.Comment {
      comment := web.Comment {
        Username: c.Product.Description,
        Avatar: c.Product.Category.Name,
        Product: c.Product.Name,
        Comment: c.Comment,
      }

      comments = append(comments, comment)
    }

    res := web.Product {
      ID: result.ID,
      Category: result.Category.Name,
      Name: result.Name,
      Description: result.Description,
      Price: result.Price,
      Available: uint32(result.QtyAvailable),
      Reserved: uint32(result.QtyReserved),
      Comment: comments,
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

  var comments []web.Comment
  for _, c := range data.Comment {
    comment := web.Comment {
      Username: c.Product.Name,
      Comment: c.Comment,
      Avatar: c.Comment,
      Product: c.Product.Name,
    }

    comments = append(comments, comment)
  }

  response := web.Product {
    ID: data.ID,
    Name: data.Name,
    Price: data.Price,
    Available: uint32(data.QtyAvailable),
    Reserved: uint32(data.QtyReserved),
    Category: data.Category.Name,
    Comment: comments,
  }

  return &response, nil
}

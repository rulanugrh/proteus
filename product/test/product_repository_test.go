package test

import (
	"testing"

	"github.com/rulanugrh/tokoku/product/internal/entity/domain"
	repomocks "github.com/rulanugrh/tokoku/product/internal/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ProductTest struct {
	suite.Suite
	repo repomocks.ProductInterface
}

func NewTestProduct() *ProductTest {
	return &ProductTest{}
}

func (product *ProductTest) TestProductCreate() {
	productResult := func (input domain.Product) *domain.Product  {
		output := &domain.Product{}
		output.Name = input.Name
		output.Price = input.Price
		output.CategoryID = input.CategoryID
		return output
	}

	product.repo.On("Create", mock.MatchedBy(func (input domain.Product) bool {
		return input.Name != "" && input.Price != 0 && input.QtyAvailable != 0
	})).Return(productResult, nil)

	data, err := product.repo.Create(domain.Product{
		Name: "MacBook 15 Pro",
		Description: "This is macbook pro",
		Price: 15000000,
		QtyAvailable: 10,
		QtyOn: 0,
		QtyReserved: 0,
		CategoryID: 1,
	})

	product.Nil(nil, err)
	product.Equal("MacBook 15 Pro", data.Name)
	product.Equal(15000000, int(data.Price))
	product.Equal(uint(1), data.CategoryID)
}

func (product *ProductTest) TestProductGetByID() {
	productResult := func (id uint) *domain.Product  {
		output := &domain.Product{}
		output.ID = id
		return output
	}

	product.repo.On("FindID", mock.MatchedBy(func (id uint) bool {
		return true
	})).Return(productResult, nil)

	data, err := product.repo.FindID(uint(1))

	product.Nil(nil, err)
	product.Equal(uint(1), data.ID)
}

func TestProduct(t *testing.T) {
	suite.Run(t, NewTestProduct())
}

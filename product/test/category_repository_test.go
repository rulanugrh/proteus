package test

import (
	"testing"

	"github.com/rulanugrh/tokoku/product/internal/entity/domain"
	repomocks "github.com/rulanugrh/tokoku/product/internal/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CategoryTest struct {
	suite.Suite
	repo repomocks.CategoryInterface
}

func NewCategoryTest() *CategoryTest {
	return &CategoryTest{}
}

func (category *CategoryTest) TestCategoryCreate() {
	categoryResult := func (input domain.Category) *domain.Category {
		output := &domain.Category{}
		output.Name = input.Name
		output.Description = input.Description
		return output
	}

	category.repo.On("Create", mock.MatchedBy(func(input domain.Category) bool {
		return input.Description != "" && input.Name != ""
	})).Return(categoryResult, nil)

	data, err := category.repo.Create(domain.Category{
		Name: "Electronic",
		Description: "This is category Electronic",
	})

	category.Nil(err)
	category.Equal("Electronic", data.Name)
	category.Equal("This is category Electronic", data.Description)
}

func (category *CategoryTest) TestCategoryFindByID() {
	categoryResult := func (id uint) *domain.Category {
		output := &domain.Category{}
		output.ID = id
		return output
	}

	category.repo.On("FindID", mock.MatchedBy(func(id uint) bool {
		return id > 0
	})).Return(categoryResult, nil)

	data, err := category.repo.FindID(uint(1))

	category.Nil(err)
	category.Equal(uint(1), data.ID)
}

func (category *CategoryTest) TestCategoryFindAll() {
	categoryResult := func () *[]domain.Category {
		output := &[]domain.Category{}
		return output
	}

	category.repo.On("FindAll", mock.Anything).Return(categoryResult, nil)

	data, err := category.repo.FindAll()

	category.Nil(err)
	category.Equal(&[]domain.Category{}, data)
}

func TestCategory(t *testing.T) {
	suite.Run(t, NewCategoryTest())
}

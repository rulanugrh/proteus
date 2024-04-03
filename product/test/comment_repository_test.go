package test

import (
	"testing"

	"github.com/rulanugrh/tokoku/product/internal/entity/domain"
	repomocks "github.com/rulanugrh/tokoku/product/internal/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CommentTest struct {
	suite.Suite
	repo repomocks.CommentInterface
}

func NewCommentTest() *CommentTest {
	return &CommentTest{}
}

func (comment *CommentTest) TestCommentCreate() {
	commentResult := func (input domain.Comment) *domain.Comment {
		output := &domain.Comment{}
		output.UserID = input.UserID
		output.ProductID = input.ProductID
		output.Rate = input.Rate
		output.RoleID = input.RoleID
		output.Comment = input.Comment
		return output
	}

	comment.repo.On("Create", mock.MatchedBy(func (input domain.Comment) bool {
		return input.UserID != 0 && input.ProductID != 0 && input.Comment != ""
	})).Return(commentResult, nil)

	data, err := comment.repo.Create(domain.Comment{
		UserID: 1,
		ProductID: 1,
		Comment: "productnya sangat bagus dan bermanfaat",
		RoleID: 3,
		Rate: 4,
	})

	comment.Nil(err)
	comment.Equal("productnya sangat bagus dan bermanfaat", data.Comment)
	comment.Equal(int8(4), data.Rate)
	comment.Equal(uint(1), data.ProductID)
	comment.Equal(uint(1), data.UserID)
	comment.Equal(uint(3), data.RoleID)
}

func (comment *CommentTest) TestCommentUID() {
	commentResult := func (id uint) *[]domain.Comment {
		output := &[]domain.Comment{}
		for _, v := range *output {
			v.UserID = id
		}

		return output
	}

	comment.repo.On("FindByUserID", mock.MatchedBy(func (id uint) bool {
		return id > 0
	})).Return(commentResult, nil)

	data, err := comment.repo.FindByUserID(1)

	for _, v := range *data {
		comment.Equal(uint(1), v.UserID)
	}
	comment.Nil(err)
}

func (comment *CommentTest) TestCommentPID() {
	commentResult := func (id uint) *[]domain.Comment {
		output := &[]domain.Comment{}
		for _, v := range *output {
			v.ProductID = id
		}

		return output
	}

	comment.repo.On("FindByProductID", mock.MatchedBy(func (id uint) bool {
		return id > 0
	})).Return(commentResult, nil)

	data, err := comment.repo.FindByProductID(1)

	for _, v := range *data {
		comment.Equal(uint(1), v.ProductID)
	}
	comment.Nil(err)
}

func TestComment(t *testing.T) {
	suite.Run(t, NewCommentTest())
}

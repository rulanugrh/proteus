package interfaces

import (
	"context"

	"github.com/ItsArul/TokoKu/entity/domain"
)

type ProductRepo interface {
	Create(product domain.Product, ctx context.Context) (domain.Product, error)
	FindById(ctx context.Context, id uint) (domain.Product, error)
	FindAll(ctx context.Context, pagination domain.PaginationProduct) []domain.Product
	Delete(ctx context.Context, id uint) error
	Update(ctx context.Context, id uint, product domain.Product) (domain.Product, error)
}

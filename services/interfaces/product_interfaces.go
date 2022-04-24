package interfaces

import (
	"context"

	"github.com/ItsArul/TokoKu/entity/domain"
)

type ProducInterfaces interface {
	Create(ctx context.Context, product domain.Product) (domain.Product, error)
	FindById(ctx context.Context, id uint) (domain.Product, error)
	FindAll(ctx context.Context) ([]domain.Product, error)
	Update(ctx context.Context, id uint, product domain.Product) (domain.Product, error)
	Delete(ctx context.Context, id uint) error
}

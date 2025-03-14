package service

import (
	"context"
	"log/slog"

	"github.com/DenisEMPS/online-shop/internal/domain"
	"github.com/DenisEMPS/online-shop/internal/filter"
	"github.com/DenisEMPS/online-shop/internal/infastructure/cache"
	"github.com/DenisEMPS/online-shop/internal/infastructure/kafka"
	"github.com/DenisEMPS/online-shop/internal/infastructure/repository"
)

type Product interface {
	Create(item *domain.CreateProduct) (int64, error)
	GetByID(id int64) (*domain.ProductDAO, error)
	GetAll(ctx context.Context, filterOptions filter.Options, sortOptions *domain.SortOptions) ([]*domain.ProductDAO, error)
}

type Auth interface {
	Register(input *domain.UserCreate) (int64, error)
	Login(input *domain.UserLogin) (string, error)
	ParseToken(token string) (int64, error)
	GenerateToken(userData *domain.UserLoginDAO) (string, error)
}

type Service struct {
	Product
	Auth
}

func NewService(repo *repository.Repository, cacheInstance cache.Cache, kafkaProd *kafka.Producer, log *slog.Logger) *Service {
	return &Service{
		Product: NewItemService(repo.Product, cacheInstance, kafkaProd, log),
		Auth:    NewAuthService(repo.Auth, kafkaProd, log),
	}
}

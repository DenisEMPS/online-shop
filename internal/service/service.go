package service

import (
	"log/slog"

	"github.com/DenisEMPS/online-shop/internal/domain"
	"github.com/DenisEMPS/online-shop/internal/infastructure/cache"
	"github.com/DenisEMPS/online-shop/internal/infastructure/repository"
)

type Item interface {
	Create(item *domain.CreateItem) (int64, error)
	GetByID(id int64) (*domain.ItemDAO, error)
}

type Auth interface {
	Register(input *domain.UserCreate) (int64, error)
	Login(input *domain.UserLogin) (string, error)
	ParseToken(token string) (int64, error)
	GenerateToken(userData *domain.UserLoginDAO) (string, error)
}

type Service struct {
	Item
	Auth
}

func NewService(repo *repository.Repository, cacheInstance cache.Cache, log *slog.Logger) *Service {
	return &Service{
		Item: NewItemService(repo.Item, cacheInstance),
		Auth: NewAuthService(repo.Auth, log),
	}
}

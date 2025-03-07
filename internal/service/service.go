package service

import (
	"github.com/DenisEMPS/online-shop/internal/domain"
	"github.com/DenisEMPS/online-shop/internal/infastructure/cache"
	"github.com/DenisEMPS/online-shop/internal/infastructure/repository"
)

type Item interface {
	Create(item *domain.CreateItem) (int64, error)
	GetByID(id int64) (*domain.ItemDAO, error)
}

type Service struct {
	Item
}

func NewService(repo *repository.Repository, cacheInstance cache.Cache) *Service {
	return &Service{
		Item: NewItemService(repo.Item, cacheInstance),
	}
}

package service

import (
	"github.com/DenisEMPS/online-shop/internal/domain"
	"github.com/DenisEMPS/online-shop/internal/infastructure/cache"
	"github.com/DenisEMPS/online-shop/internal/infastructure/repository"
)

type ItemService struct {
	repo  repository.Item
	cache cache.Cache
}

func NewItemService(repo repository.Item, cacheInstance cache.Cache) *ItemService {
	return &ItemService{repo: repo, cache: cacheInstance}
}
func (s *ItemService) Create(item *domain.CreateItem) (int64, error) {
	return s.repo.Create(item)
}
func (s *ItemService) GetByID(id int64) (*domain.ItemDAO, error) {
	if item, err := s.cache.GetItem(id); err == nil {
		return item, nil
	}

	item, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	err = s.cache.SetItem(item)

	return item, err
}

package service

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/DenisEMPS/online-shop/internal/domain"
	"github.com/DenisEMPS/online-shop/internal/infastructure/cache"
	"github.com/DenisEMPS/online-shop/internal/infastructure/repository"
)

var (
	ErrProductNotExists = errors.New("product does not exists")
)

type ItemService struct {
	repo  repository.Product
	cache cache.Cache
	log   *slog.Logger
}

func NewItemService(repo repository.Product, cacheInstance cache.Cache, log *slog.Logger) *ItemService {
	return &ItemService{repo: repo, cache: cacheInstance, log: log}
}
func (s *ItemService) Create(item *domain.CreateProduct) (int64, error) {
	const op = "item_service.create"

	log := s.log.With(
		slog.String("op", op),
		slog.String("item", item.Name),
	)

	id, err := s.repo.Create(item)
	if err != nil {
		log.Error("error to create item", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	_ = s.cache.SetItem(&domain.ProductDAO{
		ID:          id,
		Name:        item.Name,
		Description: item.Description,
		Price:       item.Price,
	})

	return id, nil
}
func (s *ItemService) GetByID(id int64) (*domain.ProductDAO, error) {
	const op = "item_service.get_by_id"

	log := s.log.With(
		slog.String("op", op),
		slog.Any("id", id),
	)

	if item, err := s.cache.GetItem(id); err == nil {
		return item, nil
	}

	item, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotExists) {
			log.Warn("product does not exists")
			return nil, fmt.Errorf("%s: %w", op, ErrProductNotExists)
		}
		log.Error("product does not exists", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_ = s.cache.SetItem(item)

	return item, err
}

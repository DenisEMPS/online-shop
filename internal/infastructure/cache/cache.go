package cache

import "github.com/DenisEMPS/online-shop/internal/domain"

type Cache interface {
	SetItem(item *domain.ItemDAO) error
	GetItem(itemID int64) (*domain.ItemDAO, error)
}

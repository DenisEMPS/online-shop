package cache

import "github.com/DenisEMPS/online-shop/internal/domain"

type Cache interface {
	SetItem(item *domain.ProductDAO) error
	GetItem(itemID int64) (*domain.ProductDAO, error)
}

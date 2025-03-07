package repository

import (
	"github.com/DenisEMPS/online-shop/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Item interface {
	Create(item *domain.CreateItem) (int64, error)
	GetByID(id int64) (*domain.ItemDAO, error)
}

type Repository struct {
	Item
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Item: NewPostgresItem(db),
	}
}

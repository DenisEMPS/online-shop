package repository

import (
	"log/slog"

	"github.com/DenisEMPS/online-shop/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Item interface {
	Create(item *domain.CreateItem) (int64, error)
	GetByID(id int64) (*domain.ItemDAO, error)
}

type Auth interface {
	Register(input *domain.UserCreate) (int64, error)
	Login(input *domain.UserLogin) (*domain.UserLoginDAO, error)
}

type Repository struct {
	Item
	Auth
}

func NewRepository(db *sqlx.DB, log *slog.Logger) *Repository {
	return &Repository{
		Item: NewPostgresItem(db),
		Auth: NewAuthPostgres(db, log),
	}
}

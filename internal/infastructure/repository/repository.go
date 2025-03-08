package repository

import (
	"log/slog"

	"github.com/DenisEMPS/online-shop/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Product interface {
	Create(item *domain.CreateProduct) (int64, error)
	GetByID(id int64) (*domain.ProductDAO, error)
}

type Auth interface {
	Register(input *domain.UserCreate) (int64, error)
	Login(input *domain.UserLogin) (*domain.UserLoginDAO, error)
}

type Repository struct {
	Product
	Auth
}

func NewRepository(db *sqlx.DB, log *slog.Logger) *Repository {
	return &Repository{
		Product: NewPostgresItem(db),
		Auth:    NewAuthPostgres(db, log),
	}
}

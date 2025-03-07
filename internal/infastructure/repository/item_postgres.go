package repository

import (
	"github.com/DenisEMPS/online-shop/internal/domain"
	"github.com/jmoiron/sqlx"
)

type ItemPostgres struct {
	db *sqlx.DB
}

func NewPostgresItem(db *sqlx.DB) *ItemPostgres {
	return &ItemPostgres{db: db}
}

func (r *ItemPostgres) Create(item *domain.CreateItem) (int64, error) {
	var id int64

	query := "INSERT INTO item (name, description, price) VALUES ($1, $2, $3) RETURNING id"

	err := r.db.QueryRow(query, item.Name, item.Description, item.Price).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
func (r *ItemPostgres) GetByID(id int64) (*domain.ItemDAO, error) {
	var item domain.ItemDAO

	query := "SELECT id, name, description, price FROM item WHERE id = $1"

	err := r.db.Get(&item, query, id)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

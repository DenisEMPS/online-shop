package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/DenisEMPS/online-shop/internal/domain"
	"github.com/jmoiron/sqlx"
)

var (
	ErrProductNotExists = errors.New("product does not exists")
)

type ProductPostgres struct {
	db *sqlx.DB
}

func NewPostgresItem(db *sqlx.DB) *ProductPostgres {
	return &ProductPostgres{db: db}
}

func (r *ProductPostgres) Create(item *domain.CreateProduct) (int64, error) {
	const op = "product_postgres.create"

	var id int64

	query := "INSERT INTO product (name, description, price) VALUES ($1, $2, $3) RETURNING id"

	err := r.db.QueryRow(query, item.Name, item.Description, item.Price).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
func (r *ProductPostgres) GetByID(id int64) (*domain.ProductDAO, error) {
	const op = "product_postgres.get_by_id"

	var item domain.ProductDAO

	query := "SELECT id, name, description, price FROM product WHERE id = $1"

	if err := r.db.QueryRow(query, id).Scan(&item.ID, &item.Name, &item.Description, &item.Price); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProductNotExists
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &item, nil
}

package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/DenisEMPS/online-shop/internal/domain"
	"github.com/DenisEMPS/online-shop/internal/filter"
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

func (r *ProductPostgres) GetAll(ctx context.Context, filterOptions filter.Options, sortOptions *domain.SortOptions) ([]*domain.ProductDAO, error) {
	const op = "product_postgres.get_all"

	var WHERE = "WHERE 1 = 1"
	var filtErr error
	var args []interface{}
	if filterOptions.IsToApply() {
		WHERE, args, filtErr = filter.BuildQuery(filterOptions)
		if filtErr != nil {
			return nil, fmt.Errorf("invalid filter options: %w", filtErr)
		}
	}

	query := fmt.Sprintf("SELECT id, name, description, price FROM product %s ORDER BY %s %s LIMIT %d", WHERE, sortOptions.SortBy, sortOptions.SortOrder, filterOptions.GetLimit())
	raws, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%v: %w with query: %v", op, err, query)
	}

	products := make([]*domain.ProductDAO, 0)
	for raws.Next() {
		var product domain.ProductDAO
		if err := raws.Scan(&product.ID, &product.Name, &product.Description, &product.Price); err != nil {
			return nil, fmt.Errorf("%v: %w with query: %v", op, err, query)
		}

		products = append(products, &product)
	}

	if len(products) == 0 {
		return nil, fmt.Errorf("%s: %w with query: %s", op, ErrProductNotExists, query)
	}

	return products, nil
}

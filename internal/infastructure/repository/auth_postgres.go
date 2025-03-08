package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/DenisEMPS/online-shop/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

type AuthPostgres struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewAuthPostgres(db *sqlx.DB, log *slog.Logger) *AuthPostgres {
	return &AuthPostgres{db: db, log: log}
}

func (r *AuthPostgres) Register(input *domain.UserCreate) (int64, error) {
	const op = "auth_postgres.register"

	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	query := "INSERT INTO adress (country, city, street) VALUES ($1, $2, $3) RETURNING id"
	var id int64

	if err := tx.QueryRow(query, input.Country, input.City, input.Street).Scan(&id); err != nil {
		errRB := tx.Rollback()
		return 0, fmt.Errorf("%s: %w\trollback_error=%w", op, err, errRB)
	}

	query = `INSERT INTO "user" (email, phone, password, first_name, second_name, adress_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	if err := tx.QueryRow(query, input.Email, input.Phone, input.PassHash, input.FirstName, input.SecondName, id).Scan(&id); err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			errRB := tx.Rollback()
			return 0, fmt.Errorf("%s: %w\trollback_error=%w", op, ErrUserExists, errRB)
		}
		errRB := tx.Rollback()
		return 0, fmt.Errorf("%s: %w\trollback_error=%w", op, err, errRB)
	}

	return id, tx.Commit()
}

func (r *AuthPostgres) Login(input *domain.UserLogin) (*domain.UserLoginDAO, error) {
	const op = "auth_postgres.login"

	query := `SELECT id, password, email FROM "user" WHERE email = $1`
	var userData domain.UserLoginDAO

	if err := r.db.QueryRow(query, input.Email).Scan(&userData.ID, &userData.PassHash, &userData.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &userData, nil
}

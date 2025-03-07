package repository

import (
	"fmt"

	"github.com/DenisEMPS/online-shop/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgres(cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.DB.Host, cfg.DB.Port, cfg.DB.Username, cfg.DB.DBname, cfg.DB.Password, cfg.DB.SSLmode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

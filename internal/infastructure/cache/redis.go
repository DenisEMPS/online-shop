package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/DenisEMPS/online-shop/internal/config"
	"github.com/DenisEMPS/online-shop/internal/domain"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	cli *redis.Client
	log *slog.Logger
}

const (
	TTLCache = 15 * time.Minute
)

func NewRedis(cfg *config.Config, log *slog.Logger) (*Redis, error) {
	redisCli := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		MaxRetries:   cfg.Redis.Maxretries,
		ReadTimeout:  cfg.Redis.ReadTimeout,
		WriteTimeout: cfg.Redis.WriteTimeout,
	})

	if err := redisCli.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis-cache: %v", err)
	}

	return &Redis{redisCli, log}, nil
}

func (r *Redis) SetItem(product *domain.ProductDAO) error {
	log := r.log.With(
		slog.Any("id", product.ID),
	)

	id := strconv.Itoa(int(product.ID))

	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	log.Info("set product to cache")
	err = r.cli.Set(context.TODO(), id, data, TTLCache).Err()
	if err != nil {
		log.Error("failed to set product in cache")
	}

	return nil
}

func (r *Redis) GetItem(productID int64) (*domain.ProductDAO, error) {
	log := r.log.With(
		slog.Any("id", productID),
	)

	var itemRes domain.ProductDAO

	id := strconv.Itoa(int(productID))

	res, err := r.cli.Get(context.TODO(), id).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, err
		}
		slog.Error("failed to get product from cache")
		return nil, err
	}

	err = json.Unmarshal([]byte(res), &itemRes)
	if err != nil {
		return nil, err
	}

	log.Info("read product from cache")

	return &itemRes, nil
}

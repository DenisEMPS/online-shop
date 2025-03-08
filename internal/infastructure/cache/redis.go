package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/DenisEMPS/online-shop/internal/config"
	"github.com/DenisEMPS/online-shop/internal/domain"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	cli *redis.Client
}

const (
	TTLCache = 15 * time.Minute
)

func NewRedis(cfg *config.Config) (*Redis, error) {
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

	return &Redis{redisCli}, nil
}

func (r *Redis) SetItem(item *domain.ItemDAO) error {
	id := strconv.Itoa(int(item.ID))

	data, err := json.Marshal(*item)
	if err != nil {
		return err
	}

	return r.cli.Set(context.TODO(), id, data, TTLCache).Err()
}

func (r *Redis) GetItem(itemID int64) (*domain.ItemDAO, error) {
	var itemOut domain.ItemDAO

	id := strconv.Itoa(int(itemID))

	res, err := r.cli.Get(context.TODO(), id).Result()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(res), &itemOut)
	if err != nil {
		return nil, err
	}

	log.Println("read item from cache")

	return &itemOut, nil
}

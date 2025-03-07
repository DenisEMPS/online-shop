package cache

import (
	"context"
	"encoding/json"
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
	rds := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		MaxRetries:   5,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	if err := rds.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &Redis{rds}, nil
}

func (r *Redis) SetItem(item *domain.ItemDAO) error {
	id := strconv.Itoa(int(item.ID))

	data, err := json.Marshal(*item)
	if err != nil {
		return err
	}

	err = r.cli.Set(context.TODO(), id, data, TTLCache).Err()
	log.Println("item saved in redis")

	return err
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

package store

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type store struct {
	rdb *redis.Client
}

type Store interface {
	GetShortURLCache(ctx context.Context, code string) string
	SetShortURLCache(ctx context.Context, code string, url string)
	RemoveShortURLCache(ctx context.Context, code string)
}

func NewStore(rdb *redis.Client) Store {
	return store{rdb}
}

func (s store) SetShortURLCache(ctx context.Context, code string, url string) {
	_, err := s.rdb.HSetEX(ctx, "urls", "EX", "86400", "FIELDS", "1", code, url).Result()
	if err != nil {
		slog.Error("erro ao salvar url na store", "err", err)
	}

}

func (s store) GetShortURLCache(ctx context.Context, code string) string {
	url, err := s.rdb.HGet(ctx, "urls", code).Result()
	if err == redis.Nil {
		return ""
	} else if err != nil {
		slog.Info("url não está na store ou já foi expirada", "err", err)
		return ""
	}

	return url
}

func (s store) RemoveShortURLCache(ctx context.Context, code string) {
	_, err := s.rdb.HDel(ctx, "urls", code).Result()
	if err != nil {
		slog.Error("erro ao remover url na store", "err", err)
	}
}

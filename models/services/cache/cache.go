package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
    redis redis.Client
    ctx context.Context
}

type ICache interface {
    Set(key string, value any) error
    Get(key string) (any, error)
    Delete(key string) error
    Clear() error
    GenerateKey(v any) (int, error)
}

func(c *Cache) Set(key string, value any) error {
    return nil
}

func(c *Cache) Get(key string) (any, error) {
    return nil, nil
}

func(c *Cache) Delete(key string) error {
    return nil
}

func(c *Cache) Clear() error {
    return nil
}

func(c *Cache) GenerateKey(v any) (int, error) {
    return 0, nil
}

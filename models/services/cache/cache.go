package cache

import (
	"context"
	"fmt"
	"os"

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

func GetRedisConnection() *redis.Client {

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

    rdb := redis.NewClient(&redis.Options{
        Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
        Password: "", 
        DB:       0, 
    })

    statusCMD := rdb.Ping(context.Background())

    if statusCMD.Err() != nil {
        panic("could not connect to redis!")
    }

    return rdb
}


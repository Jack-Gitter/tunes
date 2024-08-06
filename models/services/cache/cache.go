package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"os"
	"strconv"
	"time"

	customerrors "github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
    Redis *redis.Client
}

type ICache interface {
    Set(value any, ttl time.Duration) error
    Get(key string) (any, error)
    Delete(key string) error
    Clear() error
    GenerateKey(v any) (string, error)
}

func(c *Cache) Set(value any, ttl time.Duration) error {
    key, err := c.GenerateKey(value)

    if err != nil {
        return customerrors.WrapBasicError(err)
    }

    stringVal, err := c.TransformValueToString(value)

    if err != nil {
        return customerrors.WrapBasicError(err)
    }

    err = c.Redis.Set(context.Background(), key, stringVal, ttl).Err()

    if err != nil {
        panic(err)
    }

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

func(c *Cache) GenerateKey(v any) (string, error) {
    var bytes bytes.Buffer
    gob.NewEncoder(&bytes).Encode(v)
    return bytes.String(), nil
}

func(c *Cache) TransformValueToString(v any) (string, error) {
    var bytes bytes.Buffer
    gob.NewEncoder(&bytes).Encode(v)
    return bytes.String(), nil
}

func GetRedisConnection() *redis.Client {

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
    redisDB := os.Getenv("REDIS_DB")

    redisDBNum, err := strconv.Atoi(redisDB); 

    if err != nil {
        panic(err)
    }

    rdb := redis.NewClient(&redis.Options{
        Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
        Password: "", 
        DB:       redisDBNum, 
    })

    statusCMD := rdb.Ping(context.Background())

    if statusCMD.Err() != nil {
        panic("could not connect to redis!")
    }

    return rdb
}


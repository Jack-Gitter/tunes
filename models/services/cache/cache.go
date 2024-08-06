package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"sync"
	"time"

	customerrors "github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/Jack-Gitter/tunes/models/dtos/responses"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
    Redis *redis.Client
    Locks map[string]sync.RWMutex
}

type ICache interface {
    Set(value any, ttl time.Duration) error
    Get(key string) (any, error)
    Delete(key string) error
    Clear() error
    GenerateKey(v any) (string, error)
    LockMutex(key string) error
    UnlockMutex(key string) error
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
    switch reflect.TypeOf(v) {
        case reflect.TypeOf(responses.User{}):
            user := v.(responses.User)
            return user.SpotifyID + user.Username, nil
        case reflect.TypeOf(responses.UserIdentifer{}):
            return "", customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "caching user ids is not supported"}
        case reflect.TypeOf(responses.PostPreview{}):
            return "", customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "caching posts is not supported"}
        case reflect.TypeOf(responses.Comment{}):
            return "", customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "caching comments is not supported"}
        default: 
            return "", customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "trying to cache an unknown type!"}
    }
}

func(c *Cache) TransformValueToString(v any) (string, error) {
    var bytes bytes.Buffer
    gob.NewEncoder(&bytes).Encode(v)
    return bytes.String(), nil
}

func(c *Cache) LockMutex(key string) error {
    return nil
}

func(c *Cache) UnlockMutex(key string) error {
    return nil
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


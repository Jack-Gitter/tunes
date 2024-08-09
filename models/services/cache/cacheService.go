package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	customerrors "github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/Jack-Gitter/tunes/models/dtos/responses"
	"github.com/redis/go-redis/v9"
)

type UserCacheKey struct {
    SpotifyID string
}

type CacheService struct {
    Redis *redis.Client
    CTX context.Context
}

type ICacheService interface {
    Set(key string, value any, ttl time.Duration) error
    Get(key string) ([]byte, error)
    Delete(key string) error
    Clear() error
    GenerateKey(t reflect.Type, v any) (string, error)
}

func(c *CacheService) Set(key string, value any, ttl time.Duration) error {

    bytes, err := c.TransformValueToByteArray(value)

    if err != nil {
        return customerrors.WrapBasicError(err)
    }

    err = c.Redis.Set(c.CTX, key, bytes, ttl).Err()

    if err != nil {
        panic(err)
    }

    return nil
}

func(c *CacheService) Get(key string) ([]byte, error) {

    cmd := c.Redis.Get(c.CTX, key)

    bytes, err := cmd.Bytes()

    if err != nil {
        if errors.Is(err, redis.Nil) {
            return nil, err
        }

        return nil, &customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "cahche bad"}
    }


    return bytes, nil
}

func(c *CacheService) Delete(key string) error {
    _, err := c.Redis.Del(c.CTX, key).Result()
    if err != nil {
        return customerrors.WrapBasicError(err)
    }
    return nil
}

func(c *CacheService) Clear() error {
    _, err := c.Redis.FlushDB(c.CTX).Result()
    if err != nil {
        return customerrors.WrapBasicError(err)
    }
    return nil
}

func(c *CacheService) GenerateKey(t reflect.Type, v any) (string, error) {
    switch t {
        case reflect.TypeOf(responses.User{}): 
            user := v.(UserCacheKey)
            return user.SpotifyID, nil
        case reflect.TypeOf(responses.UserIdentifer{}):
            return "", &customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "caching user ids is not supported"}
        case reflect.TypeOf(responses.PostPreview{}):
            return "", &customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "caching posts is not supported"}
        case reflect.TypeOf(responses.Comment{}):
            return "", &customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "caching comments is not supported"}
        default: 
            return "", &customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "trying to cache an unknown type!"}
    }
}

func(c *CacheService) TransformValueToByteArray(v any) ([]byte, error) {
    var buffer bytes.Buffer

    err := gob.NewEncoder(&buffer).Encode(v)

    if err != nil {
        return []byte{}, customerrors.WrapBasicError(err)
    }

    return buffer.Bytes(), nil
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


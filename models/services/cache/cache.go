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
	"time"
	customerrors "github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/Jack-Gitter/tunes/models/dtos/responses"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
    Redis *redis.Client
    ctx context.Context
}

type ICache interface {
    Set(value any, ttl time.Duration) error
    Get(key string) ([]byte, error)
    Delete(key string) error
    Clear() error
    GenerateKey(v any) (string, error)
}

func(c *Cache) Set(key string, value any, ttl time.Duration) error {

    bytes, err := c.TransformValueToByteArray(value)

    if err != nil {
        return customerrors.WrapBasicError(err)
    }

    err = c.Redis.Set(c.ctx, key, bytes, ttl).Err()

    if err != nil {
        panic(err)
    }

    return nil
}

func(c *Cache) Get(key string) ([]byte, error) {

    cmd := c.Redis.Get(c.ctx, key)

    bytes, err := cmd.Bytes()

    if err != nil {
        panic(err)
    }

    return bytes, nil
}

func(c *Cache) Delete(key string) error {
    _, err := c.Redis.Del(c.ctx, key).Result()
    if err != nil {
        return customerrors.WrapBasicError(err)
    }
    return nil
}

func(c *Cache) Clear() error {
    _, err := c.Redis.FlushDB(c.ctx).Result()
    if err != nil {
        return customerrors.WrapBasicError(err)
    }
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

func(c *Cache) TransformValueToByteArray(v any) ([]byte, error) {
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


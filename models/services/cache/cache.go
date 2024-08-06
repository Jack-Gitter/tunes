package cache

type Cache struct {
    // is going to hold a connection to REDIS
    // implement the icache interface

}

type ICache interface {
    Set(key string, value any) error
    Get(key string) (any, error)
    Delete(key string) error
    Clear() error
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

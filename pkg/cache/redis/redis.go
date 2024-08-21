package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ofavor/ddd-go/pkg/cache"
	"github.com/ofavor/ddd-go/pkg/log"

	"github.com/redis/go-redis/v9"
)

// redisCache cache redis implementation
type redisCache struct {
	conn   *redis.Client
	prefix string
}

// NewCache create redis cache
func NewCache(addr, password string, db int32, prefix string) cache.Cache {
	log.Debug("[cache-redis] connect to ", addr)
	conn := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       int(db),
	})
	return NewCacheWithConn(conn, prefix)
}

func NewCacheWithConn(conn *redis.Client, prefix string) cache.Cache {
	return &redisCache{
		conn:   conn,
		prefix: prefix,
	}
}

func (c *redisCache) genKey(key string) string {
	return fmt.Sprintf("%s:%s", c.prefix, key)
}

// Get connection returns *redis.Client
func (c *redisCache) GetConn() interface{} {
	return c.conn
}

// Set value to cache
func (c *redisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	j, err := json.Marshal(value)
	if err != nil {
		return err
	}
	log.Debugf("[cache-redis] set %s = %s", key, j)
	err = c.conn.Set(ctx, c.genKey(key), string(j), expiration).Err()
	if err != nil {
		log.Warnf("[cache-redis] set (%s) failed: %v", key, err)
	}
	return err
}

// Get value from cache
func (c *redisCache) Get(ctx context.Context, key string, out interface{}) error {
	val, err := c.conn.Get(ctx, c.genKey(key)).Result()
	log.Debugf("[cache-redis] get %s = %s", key, val)
	if err != nil {
		if err == redis.Nil {
			//log.Warnf("[cache-redis] get (%s) failed: %v", key, err)
			return cache.ErrNil
		}
		return err
	}
	return json.Unmarshal([]byte(val), out)
}

// Delete value from cache
func (c *redisCache) Del(ctx context.Context, key ...string) error {
	nkeys := make([]string, 0, len(key))
	for _, k := range key {
		nkeys = append(nkeys, c.genKey(k))
	}
	err := c.conn.Del(ctx, nkeys...).Err()
	if err != nil {
		log.Warnf("[cache-redis] delete (%s) failed: %v", key, err)
	}
	return err
}

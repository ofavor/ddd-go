package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/ofavor/ddd-go/pkg/log"
	"github.com/ofavor/ddd-go/pkg/mutex"

	"github.com/redis/go-redis/v9"
)

type redisMutex struct {
	conn *redis.Client
}

func NewMutex(conn *redis.Client) mutex.Mutex {
	return &redisMutex{
		conn: conn,
	}
}

func (m *redisMutex) genKey(key string) string {
	return fmt.Sprintf("__locker__:%s", key)
}

func (m *redisMutex) Lock(key string, expiration time.Duration) error {
	r, err := m.conn.SetNX(context.Background(), m.genKey(key), 1, expiration).Result()
	if err != nil {
		log.Errorf("[mutex-redis] Got error while trying to lock key '%s': %v", key, err)
		return err
	}
	if !r {
		log.Errorf("[mutex-redis] Lock key '%s' failed", key)
		return mutex.ErrFail
	}
	return nil
}

func (m *redisMutex) Unlock(key string) error {
	m.conn.Del(context.Background(), m.genKey(key))
	return nil
}

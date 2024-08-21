package cache

import (
	"context"
	"errors"
	"time"
)

var ErrNil = errors.New("cache is nil")

// Cache interface
type Cache interface {
	// Get connection due to underlying implementation
	GetConn() interface{}

	// Set value to cache
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error

	// Get value from cache
	Get(ctx context.Context, key string, out interface{}) error

	// Delete value from cache
	Del(ctx context.Context, key ...string) error
}

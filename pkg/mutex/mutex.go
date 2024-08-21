package mutex

import (
	"errors"
	"time"
)

var ErrFail = errors.New("lock failed")

// Mutex interface
type Mutex interface {
	// Lock a specific key with a duration
	Lock(key string, expiration time.Duration) error

	// Unlock a specific key
	Unlock(key string) error
}

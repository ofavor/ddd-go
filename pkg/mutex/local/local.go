package local

import (
	"sync"
	"time"

	"github.com/ofavor/ddd-go/pkg/mutex"
)

type localMutex struct {
	lockers map[string]time.Time
	mutex   *sync.Mutex
}

func newLocalMutex() *localMutex {
	return &localMutex{
		lockers: make(map[string]time.Time),
		mutex:   new(sync.Mutex),
	}
}

func NewMutex() mutex.Mutex {
	return newLocalMutex()
}

func (m *localMutex) Lock(key string, expiration time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if v, ok := m.lockers[key]; ok {
		if v.After(time.Now()) {
			return mutex.ErrFail
		}
	}
	m.lockers[key] = time.Now().Add(expiration)
	return nil
}

func (m *localMutex) Unlock(key string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.lockers, key)
	return nil
}

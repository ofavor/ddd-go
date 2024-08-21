package local

import (
	"testing"
	"time"
)

func TestLockSuccess(t *testing.T) {
	m := newLocalMutex()
	key := "key.test"
	err := m.Lock(key, 0)
	if err != nil {
		t.Error("no error expected for Lock")
	}
}

func TestLockFailed(t *testing.T) {
	m := newLocalMutex()
	key := "key.test"
	err := m.Lock(key, time.Second)
	if err != nil {
		t.Error("no error expected for Lock")
	}
	err = m.Lock(key, time.Second)
	if err == nil {
		t.Error("error expected for Lock")
	}
	if err.Error() != "lock failed" {
		t.Errorf("expected error 'lock failed' but got '%s'", err.Error())
	}
}

func TestUnlock(t *testing.T) {
	m := newLocalMutex()
	key := "key.test"
	err := m.Lock(key, time.Second)
	if err != nil {
		t.Error("no error expected for Lock")
	}
	err = m.Unlock(key)
	if err != nil {
		t.Error("no error expected for Unlock")
	}
	if len(m.lockers) != 0 {
		t.Error("lockers should be empty")
	}
}

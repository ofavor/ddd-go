package redis

import (
	"errors"
	"testing"

	"github.com/go-redis/redismock/v9"
	// "time"
)

func TestLockSuccess(t *testing.T) {
	conn, mock := redismock.NewClientMock()
	red := &redisMutex{
		conn: conn,
	}
	key := "key.test"
	mock.ExpectSetNX(red.genKey(key), 1, 0).SetVal(true)
	err := red.Lock(key, 0)
	if err != nil {
		t.Error("no error expected for Lock")
	}
}

func TestLockFailed(t *testing.T) {
	conn, mock := redismock.NewClientMock()
	red := &redisMutex{
		conn: conn,
	}
	key := "key.test"
	mock.ExpectSetNX(red.genKey(key), 1, 0).SetErr(errors.New("some error"))
	err := red.Lock(key, 0)
	if err == nil {
		t.Error("error expected for Lock")
	}
	if err.Error() != "some error" {
		t.Errorf("expected error 'some error' but got '%s'", err.Error())
	}

	mock.ExpectSetNX(red.genKey(key), 1, 0).SetVal(false)

	err = red.Lock(key, 0)
	if err == nil {
		t.Error("error expected for Lock")
	}
	if err.Error() != "lock failed" {
		t.Errorf("expected error 'lock failed' but got '%s'", err.Error())
	}
}

func TestUnlockSuccess(t *testing.T) {
	conn, mock := redismock.NewClientMock()
	red := &redisMutex{
		conn: conn,
	}
	key := "key.test"
	mock.ExpectDel(red.genKey(key)).SetVal(1)
	err := red.Unlock(key)
	if err != nil {
		t.Error("no error expected for Unlock")
	}
}

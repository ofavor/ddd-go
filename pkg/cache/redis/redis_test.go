package redis

import (
	"context"
	"errors"
	"testing"

	"github.com/ofavor/ddd-go/pkg/cache"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
)

func TestSetSuccess(t *testing.T) {
	conn, mock := redismock.NewClientMock()
	red := &redisCache{
		conn:   conn,
		prefix: "myprefix",
	}
	ctx := context.Background()

	key := "key.test"
	strval := "hello"
	arrval := []string{"a", "b", "c"}
	structval := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "test",
		Age:  20,
	}

	mock.ExpectSet(red.genKey(key), "\""+strval+"\"", 0).SetVal("OK")

	err := red.Set(ctx, key, strval, 0)
	if err != nil {
		t.Error("no error expected for Set string value")
	}

	mock.ExpectSet(red.genKey(key), "[\"a\",\"b\",\"c\"]", 0).SetVal("OK")

	err = red.Set(ctx, key, arrval, 0)
	if err != nil {
		t.Error("no error expected for Set array value")
	}

	mock.ExpectSet(red.genKey(key), "{\"name\":\"test\",\"age\":20}", 0).SetVal("OK")

	err = red.Set(ctx, key, structval, 0)
	if err != nil {
		t.Error("no error expected for Set struct value")
	}
}

func TestSetFailed(t *testing.T) {
	conn, mock := redismock.NewClientMock()
	red := &redisCache{
		conn:   conn,
		prefix: "myprefix",
	}
	ctx := context.Background()

	key := "key.test"
	strval := "hello"

	mock.ExpectSet(red.genKey(key), "\""+strval+"\"", 0).SetErr(errors.New("some error"))

	err := red.Set(ctx, key, strval, 0)
	if err == nil {
		t.Error("error expected for Set string value")
	}
	if err.Error() != "some error" {
		t.Errorf("expected error 'some error' but got '%s'", err.Error())
	}
}

func TestGetSuccess(t *testing.T) {
	conn, mock := redismock.NewClientMock()
	red := &redisCache{
		conn:   conn,
		prefix: "myprefix",
	}
	ctx := context.Background()

	key := "key.test"
	strval := "hello"
	arrval := []string{"a", "b", "c"}
	structval := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "test",
		Age:  20,
	}

	mock.ExpectGet(red.genKey(key)).SetVal("\"" + strval + "\"")

	err := red.Get(ctx, key, &strval)
	if err != nil {
		t.Error("no error expected for Get string value")
	}
	if strval != "hello" {
		t.Errorf("expected 'hello' but got '%s'", strval)
	}

	mock.ExpectGet(red.genKey(key)).SetVal("[\"a\",\"b\",\"c\"]")

	err = red.Get(ctx, key, &arrval)
	if err != nil {
		t.Error("no error expected for Get array value")
	}
	if arrval[0] != "a" || arrval[1] != "b" || arrval[2] != "c" {
		t.Errorf("expected 'a,b,c' but got '%s'", arrval)
	}

	mock.ExpectGet(red.genKey(key)).SetVal("{\"name\":\"test\",\"age\":20}")

	err = red.Get(ctx, key, &structval)
	if err != nil {
		t.Error("no error expected for Get struct value")
	}
	if structval.Name != "test" || structval.Age != 20 {
		t.Errorf("expected 'test,20' but got '%s,%d'", structval.Name, structval.Age)
	}
}

func TestGetFailed(t *testing.T) {
	conn, mock := redismock.NewClientMock()
	red := &redisCache{
		conn:   conn,
		prefix: "myprefix",
	}
	ctx := context.Background()

	key := "key.test"

	mock.ExpectGet(red.genKey(key)).SetErr(errors.New("some error"))

	err := red.Get(ctx, key, &key)
	if err == nil {
		t.Error("error expected for Get string value")
	}
	if err.Error() != "some error" {
		t.Errorf("expected error 'some error' but got '%s'", err.Error())
	}

	mock.ExpectGet(red.genKey(key)).SetErr(redis.Nil)

	err = red.Get(ctx, key, &key)
	if err != cache.ErrNil {
		t.Error("no error expected for Get nil value")
	}
}

func TestDelSuccess(t *testing.T) {
	conn, mock := redismock.NewClientMock()
	red := &redisCache{
		conn:   conn,
		prefix: "myprefix",
	}
	ctx := context.Background()

	key := "key.test"

	mock.ExpectDel(red.genKey(key)).SetVal(1)

	err := red.Del(ctx, key)
	if err != nil {
		t.Error("no error expected for Del single key")
	}

	key2 := "key.test2"
	mock.ExpectDel(red.genKey(key), red.genKey(key2)).SetVal(2)

	err = red.Del(ctx, key, key2)
	if err != nil {
		t.Error("no error expected for Del multiple keys")
	}
}

func TestDelFailed(t *testing.T) {
	conn, mock := redismock.NewClientMock()
	red := &redisCache{
		conn:   conn,
		prefix: "myprefix",
	}
	ctx := context.Background()

	key := "key.test"

	mock.ExpectDel(red.genKey(key)).SetErr(errors.New("some error"))

	err := red.Del(ctx, key)
	if err == nil {
		t.Error("error expected for Del single key")
	}
	if err.Error() != "some error" {
		t.Errorf("expected error 'some error' but got '%s'", err.Error())
	}
}

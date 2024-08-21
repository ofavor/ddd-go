package kafka

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/ofavor/ddd-go/pkg/event"
)

func TestEventPubSub(t *testing.T) {
	// if true {
	// 	return
	// }
	ebus := NewEventBus([]string{"localhost:9092"}, 10, "test")
	wg := &sync.WaitGroup{}
	wg.Add(1)
	ebus.Subscribe("test", "test", func(e *event.Event) {
		fmt.Println("Event:", e)
		defer wg.Done()
	})

	ebus.Publish("test", "hello")

	wg.Wait()
}

func TestEventUnsubscribe(t *testing.T) {
	// if true {
	// 	return
	// }
	h := func(e *event.Event) {
		fmt.Println("Event:", e)
	}
	ebus := NewEventBus([]string{"localhost:9092"}, 10, "test")
	ebus.Subscribe("test", "test", h)
	// ebus.Subscribe("test", "test1", h)
	time.Sleep(time.Second * 1)
	// ebus.Publish("test", "hello")
	// time.Sleep(time.Second * 5)
	ebus.Unsubscribe("test", "test", h)
	// ebus.Unsubscribe("test", "test1", h)
	for {
		time.Sleep(time.Second)
	}
}

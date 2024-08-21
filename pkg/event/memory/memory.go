package memory

import (
	"reflect"
	"runtime/debug"
	"sync"

	"github.com/ofavor/ddd-go/pkg/event"
	"github.com/ofavor/ddd-go/pkg/log"
)

// Memory event bus implementation.
type eventHandlerWrapper struct {
	bus        *memoryEventBus
	moduleName string
	callback   reflect.Value
	fn         event.EventHandler
	events     chan *event.Event
}

func (h *eventHandlerWrapper) start() {
	go func() {
		for e := range h.events {
			h.handle(e)
		}
	}()
}

func (h *eventHandlerWrapper) handle(e *event.Event) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("[event-mem] Got error while handling event:", e)
			log.Error(err, string(debug.Stack()))
		}
	}()
	h.fn(e)
}

func (h *eventHandlerWrapper) stop() {
	if h.events != nil {
		close(h.events)
	}
}

type memoryEventBus struct {
	handlers map[string][]*eventHandlerWrapper
	lock     *sync.RWMutex
	events   chan *event.Event
}

const consumerBufferSize = 100

// NewEventBus create memory event bus
func NewEventBus(bufferSize int64) event.EventBus {
	bus := &memoryEventBus{
		handlers: map[string][]*eventHandlerWrapper{},
		lock:     new(sync.RWMutex),
		events:   make(chan *event.Event, bufferSize),
	}
	go bus.run() // start main goroutine
	return bus
}

func (b *memoryEventBus) run() {
	for e := range b.events {
		log.Debug("[event-mem] Bus got event:", e)
		b.handleEvent(e)
	}
}

func (b *memoryEventBus) handleEvent(e *event.Event) {
	b.lock.RLock()
	defer b.lock.RUnlock()
	hh, ok := b.handlers[e.Meta().Type]
	if ok {
		for _, k := range hh {
			k.events <- e
		}
	}
}

func (b *memoryEventBus) Publish(t string, payload interface{}) error {
	e, _ := event.NewEvent(t, payload)
	b.events <- e
	return nil
}

func (b *memoryEventBus) Subscribe(t string, name string, h event.EventHandler) error {
	log.Debug("[event-mem] Subscribe event: ", t)
	b.lock.Lock()
	defer b.lock.Unlock()
	hh, ok := b.handlers[t]
	if !ok {
		hh = []*eventHandlerWrapper{}
	}
	eh := &eventHandlerWrapper{bus: b, moduleName: name, callback: reflect.ValueOf(h), fn: h, events: make(chan *event.Event, consumerBufferSize)}
	eh.start()
	hh = append(hh, eh)
	b.handlers[t] = hh
	return nil
}

func (b *memoryEventBus) Unsubscribe(t string, name string, h event.EventHandler) error {
	log.Debug("[event-mem] Unsubscribe event: ", t)
	b.lock.Lock()
	defer b.lock.Unlock()
	c := reflect.ValueOf(h)
	hh, ok := b.handlers[t]
	if ok {
		for i, k := range hh {
			if k.moduleName == name && k.callback.Pointer() == c.Pointer() {
				hh = append(hh[:i], hh[i+1:]...)
				k.stop()
				break
			}
		}
		if len(hh) == 0 {
			delete(b.handlers, t)
		}
	}
	return nil
}

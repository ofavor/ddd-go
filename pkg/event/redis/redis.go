package redis

import (
	"context"
	"fmt"
	"reflect"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"github.com/ofavor/ddd-go/pkg/event"
	"github.com/ofavor/ddd-go/pkg/log"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type eventHandlerWrapper struct {
	bus        *redisEventBus
	moduleName string
	callback   reflect.Value
	fn         event.EventHandler
	events     chan *event.Event
}

func (h *eventHandlerWrapper) start() {
	for e := range h.events {
		h.handle(e)
	}
}

func (h *eventHandlerWrapper) handle(e *event.Event) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("[event-redis] Got error while handling event:", e)
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

const consumerBufferSize int64 = 1000

type eventConsumer struct {
	bus       *redisEventBus
	eventType string
	handlers  []*eventHandlerWrapper
	cancel    context.CancelFunc
}

func (c *eventConsumer) start() {
	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel
	streamKey := c.bus.genStreamKey(c.eventType)
	for {
		log.Debugf("[event-redis] Trying to create consumer group: %s %s", streamKey, c.bus.group)
		err := c.bus.conn.XGroupCreate(ctx, streamKey, c.bus.group, "0").Err()
		if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
			log.Warnf("[event-redis] Got error while creating consumer group: %v", err)
			time.Sleep(time.Second * 10)
		} else {
			break
		}
	}
	cid := uuid.NewString()
	for {
		log.Debugf("[event-redis] Trying to read streams from redis: %s %s %s", streamKey, c.bus.group, cid)
		// TODO context.cancel cannot stop XReadGroup block
		if streams, err := c.bus.conn.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    c.bus.group,
			Consumer: cid,
			Count:    1,
			Streams:  []string{streamKey, ">"},
			Block:    0,
		}).Result(); err != nil {
			if err == context.Canceled { // context.cancel cannot stop XReadGroup block, this error might be not received ever
				log.Infof("[event-redis] Consumer group %s canceled", c.bus.group)
				break
			}
			log.Warnf("[event-redis] Got error while consuming events: %v", err)
		} else {
			if len(c.handlers) == 0 { // no handlers, break the loop
				break
			}
			// log.Debug("[event-redis] Got events from: ", c.subject)
			for _, stream := range streams {
				for _, msg := range stream.Messages {
					// log.Debug("[event-redis] Consume event: ", msg.ID)
					if err := c.bus.conn.XAck(ctx, c.eventType, c.bus.group, msg.ID).Err(); err != nil {
						log.Warnf("[event-redis] Got error while acking event: %v", err)
					}
					id := msg.Values["id"].(string)
					tm, _ := strconv.ParseInt(msg.Values["time"].(string), 10, 64)
					typee := msg.Values["type"].(string)
					payload := msg.Values["payload"].(string)
					e, err := event.LoadEvent(id, tm, typee, payload)
					if err != nil {
						log.Warnf("[event-redis] Got error while loading event: %v", err)
						continue
					}
					log.Debug("[event-redis] Dispatch event to handlers: ", e)
					for _, h := range c.handlers {
						h.events <- e
					}
				}
			}
		}
	}
}

func (c *eventConsumer) stop() {
	if c.cancel != nil {
		log.Debug("[event-redis] Cancel event consumer: ", c.eventType)
		c.cancel()
	}
}

type redisEventBus struct {
	conn       *redis.Client
	bufferSize int64
	group      string
	consumers  map[string]*eventConsumer
	lock       *sync.RWMutex
}

func NewEventBus(addr, password string, db int32, bufferSize int64, group string) event.EventBus {
	conn := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       int(db),
	})
	return NewEventBusWithConn(conn, bufferSize, group)
}

func NewEventBusWithConn(conn *redis.Client, bufferSize int64, group string) event.EventBus {
	bus := &redisEventBus{
		conn:       conn,
		bufferSize: bufferSize,
		group:      group,
		consumers:  map[string]*eventConsumer{},
		lock:       new(sync.RWMutex),
	}
	return bus
}

func (b *redisEventBus) genStreamKey(t string) string {
	return fmt.Sprintf("__event__:%s", t)
}

// Publish implements event.EventBus.
func (b *redisEventBus) Publish(t string, payload interface{}) error {
	e, _ := event.NewEvent(t, payload)
	log.Debug("[event-redis] Publish event: ", e)
	if err := b.conn.XAdd(context.Background(), &redis.XAddArgs{
		Stream: b.genStreamKey(e.Meta().Type),
		MaxLen: b.bufferSize,
		Values: map[string]interface{}{
			"id":      e.Id().String(),
			"type":    e.Meta().Type,
			"time":    e.Meta().Time.UnixNano(),
			"payload": string(e.Payload()),
		},
	}).Err(); err != nil {
		log.Warnf("[event-redis] Got error while publishing event: %v", err)
		return err
	}
	return nil
}

// Subscribe implements event.EventBus.
func (b *redisEventBus) Subscribe(t string, name string, h event.EventHandler) error {
	log.Debug("[event-redis] Subscribe event: ", t)
	b.lock.Lock()
	defer b.lock.Unlock()
	c, ok := b.consumers[t]
	if !ok {
		log.Debugf("[event-redis] Create and start consumer for event: %s", t)
		c = &eventConsumer{
			bus:       b,
			eventType: t,
			handlers:  []*eventHandlerWrapper{},
		}
		b.consumers[t] = c
		go c.start()
	}
	eh := &eventHandlerWrapper{bus: b, moduleName: name, callback: reflect.ValueOf(h), fn: h, events: make(chan *event.Event, consumerBufferSize)}
	go eh.start()
	c.handlers = append(c.handlers, eh)
	return nil
}

// Unsubscribe implements event.EventBus.
func (b *redisEventBus) Unsubscribe(t string, name string, h event.EventHandler) error {
	log.Debug("[event-redis] Unsubscribe event: ", t)
	b.lock.Lock()
	defer b.lock.Unlock()
	rh := reflect.ValueOf(h)
	c, ok := b.consumers[t]
	if ok {
		for i, k := range c.handlers {
			if k.moduleName == name && k.callback.Pointer() == rh.Pointer() {
				c.handlers = append(c.handlers[:i], c.handlers[i+1:]...)
				k.stop()
				break
			}
		}
		if len(c.handlers) == 0 {
			log.Debugf("[event-redis] Delete and stop consumer for event: %s", t)
			delete(b.consumers, t)
			c.stop()
		}
	}
	return nil
}

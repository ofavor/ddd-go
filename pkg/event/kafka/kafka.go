package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"github.com/ofavor/ddd-go/pkg/event"
	"github.com/ofavor/ddd-go/pkg/log"

	"github.com/segmentio/kafka-go"
)

type eventHandlerWrapper struct {
	bus        *kafkaEventBus
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
			log.Error("[event-kafka] Got error while handling event:", e)
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
	bus       *kafkaEventBus
	eventType string
	handlers  []*eventHandlerWrapper
	cancel    context.CancelFunc
}

func (c *eventConsumer) prepareTopic() error {
	conn, err := kafka.Dial("tcp", c.bus.brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()
	ctl, err := conn.Controller()
	if err != nil {
		return err
	}
	cconn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", ctl.Host, ctl.Port))
	if err != nil {
		return err
	}
	defer cconn.Close()
	err = cconn.CreateTopics(kafka.TopicConfig{
		Topic: c.bus.genTopicKey(c.eventType),
	})
	return err
}

func (c *eventConsumer) start() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: c.bus.brokers,
		Topic:   c.bus.genTopicKey(c.eventType),
		GroupID: c.bus.group,
	})
	defer reader.Close()
	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel

	if err := c.prepareTopic(); err != nil {
		log.Warnf("[event-kafka] Got error while preparing topic: %v", err)
		// time.Sleep(time.Second * 10)
		// continue
		return
	}

	for {
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			if err.Error() == "fetching message: context canceled" { // context is canceled
				log.Info("[event-kafka] Context canceled")
				break
			}
			log.Warnf("[event-kafka] Got error while reading message: %v", err)
			time.Sleep(time.Second * 10)
			continue
		}
		val := map[string]interface{}{}
		err = json.Unmarshal(m.Value, &val)
		if err != nil {
			log.Warnf("[event-kafka] Got error while unmarshalling message: %v", err)
			continue
		}
		id := val["id"].(string)
		tm, _ := strconv.ParseInt(val["time"].(string), 10, 64)
		tp := val["type"].(string)
		pl := val["payload"].(string)
		e, err := event.LoadEvent(id, tm, tp, pl)
		if err != nil {
			log.Warnf("[event-kafka] Got error while loading event: %v", err)
			continue
		}
		for _, h := range c.handlers {
			h.events <- e
		}
	}
}

func (c *eventConsumer) stop() {
	if c.cancel != nil {
		log.Debug("[event-kafka] Cancel event consumer: ", c.eventType)
		c.cancel()
	}
}

type kafkaEventBus struct {
	brokers    []string
	bufferSize int64
	group      string
	consumers  map[string]*eventConsumer
	lock       *sync.RWMutex
}

func NewEventBus(brokers []string, bufferSize int64, group string) event.EventBus {
	bus := &kafkaEventBus{
		brokers:    brokers,
		bufferSize: bufferSize,
		group:      group,
		consumers:  map[string]*eventConsumer{},
		lock:       new(sync.RWMutex),
	}
	return bus
}

func (b *kafkaEventBus) genTopicKey(t string) string {
	return fmt.Sprintf("__event__.%s", t)
}

// Publish implements event.EventBus.
func (b *kafkaEventBus) Publish(t string, payload interface{}) error {
	e, _ := event.NewEvent(t, payload)
	log.Debug("[event-kafka] Publish event: ", e)
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: b.brokers,
		Topic:   b.genTopicKey(t),
	})
	ev, err := json.Marshal(map[string]interface{}{
		"id":      e.Id().String(),
		"time":    fmt.Sprintf("%d", e.Meta().Time.UnixNano()),
		"type":    e.Meta().Type,
		"payload": e.Payload(),
	})
	if err != nil {
		return err
	}
	if err := writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(e.Id().String()),
		Value: ev,
	}); err != nil {
		log.Warnf("[event-kafka] Got error while publishing event: %v", err)
		return err
	}
	return nil
}

// Subscribe implements event.EventBus.
func (b *kafkaEventBus) Subscribe(t string, name string, h event.EventHandler) error {
	log.Debug("[event-kafka] Subscribe event: ", t)
	b.lock.Lock()
	defer b.lock.Unlock()
	c, ok := b.consumers[t]
	if !ok {
		log.Debugf("[event-kafka] Create and start consumer for event: %s", t)
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
func (b *kafkaEventBus) Unsubscribe(t string, name string, h event.EventHandler) error {
	log.Debug("[event-kafka] Unsubscribe event: ", t)
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
			log.Debugf("[event-kafka] Delete and stop consumer for event: %s", t)
			delete(b.consumers, t)
			c.stop()
		}
	}
	return nil
}

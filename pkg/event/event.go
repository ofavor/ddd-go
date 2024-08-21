package event

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Event meta
type Meta struct {
	Type string
	Time time.Time
}

// Event
type Event struct {
	id      uuid.UUID
	meta    Meta
	payload json.RawMessage
}

// Create event
func NewEvent(t string, data interface{}) (*Event, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	meta := Meta{
		Type: t,
		Time: time.Now(),
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return &Event{
		id:      id,
		meta:    meta,
		payload: payload,
	}, nil
}

func LoadEvent(id string, tm int64, t string, data string) (*Event, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	meta := Meta{
		Type: t,
		Time: time.Unix(0, tm),
	}
	return &Event{
		id:      uid,
		meta:    meta,
		payload: []byte(data),
	}, nil
}

// Get event Id
func (e *Event) Id() uuid.UUID {
	return e.id
}

// Get event meta
func (e *Event) Meta() Meta {
	return e.meta
}

// Get event payload
func (e *Event) Payload() json.RawMessage {
	return e.payload
}

// Get event string
func (e *Event) String() string {
	return fmt.Sprintf("event{id=%s type=%s time=%s payload=%s}", e.id, e.meta.Type, e.meta.Time, string(e.payload))
}

// EventHandler consume event
type EventHandler func(e *Event /*, errReporter EventErrReporter*/)

// Bus bus
type EventBus interface {
	// Publish event by specifying event type and payload
	Publish(t string, payload interface{}) error

	// Subscribe event handler
	Subscribe(t string, name string, h EventHandler) error

	// Unsubscribe event handler
	Unsubscribe(t string, name string, h EventHandler) error
}

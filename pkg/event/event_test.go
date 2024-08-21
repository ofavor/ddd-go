package event

import (
	"testing"
)

func TestEventTime(t *testing.T) {
	e, _ := NewEvent("test", "hello")
	e1, _ := LoadEvent(e.id.String(), e.meta.Time.UnixNano(), "test", "hello")
	// fmt.Println(">>>>>>>>>>", e.meta.Time, "vs", e1.meta.Time)
	if e.meta.Time.UnixNano() != e1.meta.Time.UnixNano() {
		t.Error("event time error")
	}
}

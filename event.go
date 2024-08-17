package event

import (
	"sync"
)

var closedchan = make(chan struct{})

func init() {
	close(closedchan)
}

// Event represents a one-time event that may occur in the future.
type Event struct {
	// trigger  atomic.Uint32
	done     chan struct{}
	initOnce sync.Once
}

// Trigger causes e to complete.  It is safe to call multiple times, and
// concurrently.  It returns true if this call does cause the event to fire.
func (e *Event) Trigger() bool {
	flag := false
	e.initOnce.Do(func() {
		e.done = closedchan
		flag = true
	})
	if flag {
		return true
	} else if e.HasTrigger() {
		return false
	} else {
		close(e.done)
		return true
	}
}

// Done returns a channel that will be closed when Fire is called.
func (e *Event) Done() <-chan struct{} {
	e.initOnce.Do(func() {
		e.done = make(chan struct{})
	})
	return e.done
}

// HasTrigger returns true if Trigger has been called.
func (e *Event) HasTrigger() bool {
	select {
	case <-e.Done():
		return true
	default:
		return false
	}
}

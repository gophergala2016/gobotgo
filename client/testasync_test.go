package client

import (
	"testing"
	"time"
)

// asyncronous execution validator
// Should factor this into a separate package as testasync.Order
// Also add a "wait for event > N" function, Wait(order int)
type validate struct {
	ch    chan int
	count int
	t     *testing.T
}

func validator(t *testing.T) *validate {
	return &validate{
		ch: make(chan int, 1000), // we're essentially using the validator as a syncronous slice
		t:  t,
	}
}

func (v *validate) Verify(d time.Duration) {
	position := 0
	count := 0
	timeout := time.After(d)
	for {
		select {
		case order := <-v.ch:
			v.t.Logf("Received event %d", order)
			if order < position {
				v.t.Errorf("Received event %d after %d", order, position)
			}
			position = order
			count++
			if count == v.count {
				v.count = 0
				return
			}
		case <-timeout:
			v.t.Errorf("verify failed after %s", d)
			v.t.Fatalf("only finished %d of %d expected events", count, v.count)
			return
		}
	}
}

func (v *validate) Before(order int, f func()) {
	v.count++
	go func() {
		v.ch <- order
		f()
	}()
}

func (v *validate) After(order int, f func()) {
	v.count++
	go func() {
		f()
		v.ch <- order
	}()
}

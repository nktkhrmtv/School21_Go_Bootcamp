package main

import (
	"testing"
	"time"
)

func TestMultiplex(t *testing.T) {
	ch1 := make(chan interface{})
	ch2 := make(chan interface{})
	ch3 := make(chan interface{})

	out := multiplex(ch1, ch2, ch3)

	go func() {
		defer close(ch1)
		defer close(ch2)
		defer close(ch3)
		ch1 <- "ch1-1"
		ch2 <- "ch2-1"
		ch3 <- "ch3-1"
		ch1 <- "ch1-2"
		ch2 <- "ch2-2"
	}()

	expected := map[string]bool{
		"ch1-1": false,
		"ch1-2": false,
		"ch2-1": false,
		"ch2-2": false,
		"ch3-1": false,
	}

	for i := 0; i < 5; i++ {
		select {
		case val := <-out:
			if _, ok := expected[val.(string)]; !ok {
				t.Errorf("Unexpected value: %v", val)
			}
			expected[val.(string)] = true
		case <-time.After(100 * time.Millisecond):
			t.Error("Timeout waiting for value")
		}
	}

	for k, v := range expected {
		if !v {
			t.Errorf("Expected value not received: %s", k)
		}
	}
}
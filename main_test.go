package main

import (
	"testing"
	"time"
)

func TestRunSession(t *testing.T) {
	s := runSession(time.Second*5, "test", true)

	if s.Duration != 0 {
		t.Errorf("got %d, want %d", s.Duration, 5)
	}
	if s.Category != "test" {
		t.Errorf("got %s, want %s", s.Category, "test")
	}
}

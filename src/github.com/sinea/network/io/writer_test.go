package io

import (
	"testing"
	"time"
)

func TestIOWrite(t *testing.T) {
	sw := NewWriter(nil)
	sw.Write([]byte("hello"))
	sw.Write([]byte("world"))

	time.Sleep(time.Hour)
}
package server

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	t.Run("Hit", func(t *testing.T) {
		c := NewCache(5, 5*time.Second)
		c.Add("test", [][][]string{{{"test"}}})

		_, ok := c.Get("test")

		if !ok {
			t.Errorf("expected item to exist")
		}
	})

	t.Run("Miss", func(t *testing.T) {
		c := NewCache(5, 5*time.Second)

		_, ok := c.Get("test")

		if ok {
			t.Errorf("expected item to not exist")
		}
	})

	t.Run("Expiration", func(t *testing.T) {
		c := NewCache(5, 500*time.Millisecond)
		c.Add("test", [][][]string{{{"test"}}})

		time.Sleep(1 * time.Second)

		_, ok := c.Get("test")

		if ok {
			t.Errorf("expected item to not exist")
		}
	})
}

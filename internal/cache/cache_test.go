package cache

import (
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func TestGet(t *testing.T) {
	t.Run("Hit", func(t *testing.T) {
		c := New(5, 5*time.Second, 5*time.Second)
		c.Set("test", new(goquery.Selection))

		_, ok := c.Get("test")

		if !ok {
			t.Errorf("expected item to exist")
		}
	})

	t.Run("Miss", func(t *testing.T) {
		c := New(5, 5*time.Second, 5*time.Second)

		_, ok := c.Get("test")

		if ok {
			t.Errorf("expected item to not exist")
		}
	})

	t.Run("Expiration", func(t *testing.T) {
		c := New(5, 500*time.Millisecond, 500*time.Millisecond)
		c.Set("test", new(goquery.Selection))

		time.Sleep(1 * time.Second)

		_, ok := c.Get("test")

		if ok {
			t.Errorf("expected item to not exist")
		}
		c.Print()
	})
}

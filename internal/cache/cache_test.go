package cache

import (
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func TestGet(t *testing.T) {
	t.Run("Hit", func(t *testing.T) {
		c := NewCache(5, 5*time.Second, 5*time.Second)
		c.Set("test", new(goquery.Selection))

		_, ok := c.Get("test")

		if !ok {
			t.Errorf("expected item to exist")
		}

		if len(c.elements) != 1 || c.list.Len() != 1 {
			t.Errorf("expected one item in the cache")
		}
	})

	t.Run("Miss", func(t *testing.T) {
		c := NewCache(5, 5*time.Second, 5*time.Second)

		_, ok := c.Get("test")

		if ok {
			t.Errorf("expected item to not exist")
		}
	})

	t.Run("Expiration", func(t *testing.T) {
		c := NewCache(5, 500*time.Millisecond, 500*time.Millisecond)
		c.Set("test", new(goquery.Selection))

		time.Sleep(1 * time.Second)

		_, ok := c.Get("test")

		if ok {
			t.Errorf("expected item to not exist")
		}

		if len(c.elements) != 0 || c.list.Len() != 0 {
			t.Errorf("expected no items in the cache")
		}
	})

	t.Run("Capacity", func(t *testing.T) {
		c := NewCache(2, 5*time.Second, 5*time.Second)
		c.Set("one", new(goquery.Selection))
		c.Set("two", new(goquery.Selection))
		c.Set("three", new(goquery.Selection))

		_, ok := c.Get("two")

		if !ok {
			t.Errorf("expected two item to exist")
		}

		_, ok = c.Get("three")

		if !ok {
			t.Errorf("expected three item to exist")
		}

		if len(c.elements) != 2 || c.list.Len() != 2 {
			t.Errorf("expected two items in the cache")
		}
	})
}

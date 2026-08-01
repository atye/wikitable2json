package server

import (
	"encoding/json"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	t.Run("Hit", func(t *testing.T) {
		c := NewCache(5, 5*time.Second)
		key := marshalCacheKey(t, cacheKey{Page: "test"})
		c.Add(key, [][][]string{{{"test"}}})

		_, ok := c.Get(key)

		if !ok {
			t.Errorf("expected item to exist")
		}
	})

	t.Run("Miss", func(t *testing.T) {
		c := NewCache(5, 5*time.Second)
		key := marshalCacheKey(t, cacheKey{Page: "test"})

		_, ok := c.Get(key)

		if ok {
			t.Errorf("expected item to not exist")
		}
	})

	t.Run("Expiration", func(t *testing.T) {
		c := NewCache(5, 500*time.Millisecond)
		key := marshalCacheKey(t, cacheKey{Page: "test"})
		c.Add(key, [][][]string{{{"test"}}})

		time.Sleep(1 * time.Second)

		_, ok := c.Get(key)

		if ok {
			t.Errorf("expected item to not exist")
		}
	})
}
func marshalCacheKey(t *testing.T, key cacheKey) string {
	t.Helper()

	b, err := json.Marshal(key)
	if err != nil {
		t.Fatalf("marshal cache key: %v", err)
	}
	return string(b)
}

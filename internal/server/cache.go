package server

import (
	"log"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
)

type Cache struct {
	lru *expirable.LRU[string, any]
}

func NewCache(size int, expiration time.Duration) *Cache {
	return &Cache{
		lru: expirable.NewLRU[string, any](size, onEvict, expiration),
	}
}

func (c *Cache) Get(key string) (any, bool) {
	return c.lru.Get(key)
}

func (c *Cache) Add(key string, value any) bool {
	log.Printf("cache: added %v\n", key)
	return c.lru.Add(key, value)
}

func onEvict[K comparable, V any](key K, value V) {
	log.Printf("cache: evicted %v\n", key)
}

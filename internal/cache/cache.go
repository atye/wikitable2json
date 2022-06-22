package cache

import (
	"container/list"
	"log"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type item struct {
	key       string
	selection *goquery.Selection
	expiresAt time.Time
}

type Cache struct {
	mu         sync.Mutex
	capacity   int
	list       *list.List               //DoublyLinkedList for backing the cache value.
	elements   map[string]*list.Element //Map to store list pointer of cache mapped to key
	expiration time.Duration
	purgeEvery time.Duration
	log        *log.Logger
}

func New(capacity int, itemExpiration time.Duration, purgeEvery time.Duration) *Cache {
	c := &Cache{
		mu:         sync.Mutex{},
		capacity:   capacity,
		expiration: itemExpiration,
		purgeEvery: purgeEvery,
		list:       new(list.List),
		elements:   make(map[string]*list.Element),
		log:        log.Default(),
	}

	go func() {
		ticker := time.NewTicker(purgeEvery)
		for range ticker.C {
			c.mu.Lock()
			for key := range c.elements {
				if time.Now().After(c.elements[key].Value.(*list.Element).Value.(item).expiresAt) {
					c.log.Printf("remove from cache: %s", key)
					c.list.Remove(c.elements[key])
					delete(c.elements, key)
				}
			}
			c.mu.Unlock()
		}
	}()

	return c
}

func (c *Cache) Get(key string) (*goquery.Selection, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, ok := c.elements[key]; ok {
		c.log.Printf("cache hit: %s", key)
		item := element.Value.(*list.Element).Value.(item)
		item.expiresAt = time.Now().Add(c.expiration)
		element.Value = item
		c.list.MoveToFront(element)
		return item.selection, true
	}
	return nil, false
}

func (c *Cache) Set(key string, s *goquery.Selection) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, ok := c.elements[key]; ok {
		element.Value = item{key: key, selection: s, expiresAt: time.Now().Add(c.expiration)}
		c.list.MoveToFront(element)
	} else {
		if c.list.Len() == c.capacity {
			lruElement := c.list.Back()
			lruKey := lruElement.Value.(*list.Element).Value.(item).key
			c.list.Remove(lruElement)
			delete(c.elements, lruKey)
		}

		element := &list.Element{
			Value: item{
				key:       key,
				selection: s,
				expiresAt: time.Now().Add(c.expiration),
			},
		}

		pointer := c.list.PushFront(element)
		c.elements[key] = pointer
	}

	c.log.Printf("add to cache: %s", key)
}

func (c *Cache) Print() {
	c.log.Println(c.list.Len())
	c.log.Println(len(c.elements))
}

package client

import (
	"container/list"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type item struct {
	key       string
	selection *goquery.Selection
	expiresAt time.Time
}

type cache struct {
	mu             sync.Mutex
	capacity       int
	list           *list.List
	elements       map[string]*list.Element
	itemExpiration time.Duration
	purgeEvery     time.Duration
}

func newCache(capacity int, itemExpiration time.Duration, purgeEvery time.Duration) *cache {
	c := &cache{
		mu:             sync.Mutex{},
		capacity:       capacity,
		itemExpiration: itemExpiration,
		purgeEvery:     purgeEvery,
		list:           new(list.List),
		elements:       make(map[string]*list.Element),
	}

	go func() {
		ticker := time.NewTicker(purgeEvery)
		for range ticker.C {
			c.mu.Lock()
			for key := range c.elements {
				if time.Now().After(c.elements[key].Value.(*list.Element).Value.(item).expiresAt) {
					c.list.Remove(c.elements[key])
					delete(c.elements, key)
				}
			}
			c.mu.Unlock()
		}
	}()

	return c
}

func (c *cache) Get(key string) (*goquery.Selection, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, ok := c.elements[key]; ok {
		c.list.MoveToFront(element)
		return element.Value.(*list.Element).Value.(item).selection, true
	}
	return nil, false
}

func (c *cache) Set(key string, s *goquery.Selection) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, ok := c.elements[key]; ok {
		element.Value.(*list.Element).Value = item{key: key, selection: s, expiresAt: element.Value.(*list.Element).Value.(item).expiresAt}
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
				expiresAt: time.Now().Add(c.itemExpiration),
			},
		}

		pointer := c.list.PushFront(element)
		c.elements[key] = pointer
	}
}

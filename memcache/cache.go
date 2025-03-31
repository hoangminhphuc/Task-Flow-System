package memcache

import "sync"

type Caching interface {
    Write(k string, value interface{})
    Read(k string) interface{}
}

type caching struct {
    store  map[string]interface{}
    locker *sync.RWMutex
}

func NewCaching() *caching {
    return &caching{
        store:  make(map[string]interface{}),
        locker: new(sync.RWMutex),
    }
}

func (c *caching) Write(k string, value interface{}) {
	c.locker.Lock()
	defer c.locker.Unlock()
	c.store[k] = value
}

func (c *caching) Read(k string) interface{} {
	c.locker.RLock()
	defer c.locker.RUnlock()
	return c.store[k]
}
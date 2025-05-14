package cache

import (
	"sync"
	"time"
)

type item struct {
	value      interface{}
	expiration int64
}

type Cache struct {
	data  map[string]item
	mutex sync.RWMutex
	ttl   time.Duration
}

var (
	defaultCache *Cache
	once         sync.Once
)

func getCache() *Cache {
	once.Do(func() {
		defaultCache = &Cache{
			data: make(map[string]item),
			ttl:  10 * time.Minute,
		}
	})
	return defaultCache
}

func Set(key string, value interface{}) {
	cache := getCache()

	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.data[key] = item{
		value:      value,
		expiration: time.Now().Add(cache.ttl).UnixNano(),
	}
}

func Get(key string) (interface{}, bool) {
	cache := getCache()

	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	it, found := cache.data[key]
	if !found || time.Now().UnixNano() > it.expiration {
		return nil, false
	}
	return it.value, true
}

func Delete(key string) {
	cache := getCache()

	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	delete(cache.data, key)
}

package storage

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	cacheOnce     sync.Once
	cacheInstance *cache.Cache
)

const (
	DefaultExpiration = cache.DefaultExpiration
)

func Cache() *cache.Cache {
	if cacheInstance == nil {
		cacheOnce.Do(func() {
			cacheInstance = cache.New(5*time.Minute, 10*time.Minute)
		})
	}

	return cacheInstance
}

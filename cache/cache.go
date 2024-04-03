package cache

import (
	"log"
	"math/rand/v2"
	"time"

	"github.com/jellydator/ttlcache/v3"
)

func Run() {
	cacheLoop(3*time.Second, 1*time.Second, 10)
}

func cacheLoop(cacheTTL time.Duration, loopSleep time.Duration, loops int) {
	cache := ttlcache.New[string, int](
		ttlcache.WithTTL[string, int](cacheTTL),
		// Unless this is disabled, Get() also extends/touches an item's expiration timestamp on successful retrieval.
		ttlcache.WithDisableTouchOnHit[string, int](),
	)

	go cache.Start() // starts automatic expired item deletion
	for i := 0; i < loops; i++ {
		key := "rand"
		if cache.Has(key) {
			item := cache.Get(key) // If the item is not found, a nil value is returned.
			log.Printf("Got key %s from cache val %d expires at %v\n", key, item.Value(), item.ExpiresAt().Format("15:04:05 MST"))
		} else {
			r := rand.IntN(100)
			log.Printf("Got key %s set new val %d with default cacheTTL\n", key, r)
			cache.Set(key, r, ttlcache.DefaultTTL)
		}
		// cache.DeleteExpired() // no need if cache.Start() runs as separate goroutine
		time.Sleep(loopSleep)
	}
	cache.Stop()
}

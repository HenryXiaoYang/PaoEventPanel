package handler

import (
	"sync"
	"time"
)

type cacheEntry struct {
	data      any
	expiresAt time.Time
}

var (
	rankingCache = make(map[string]cacheEntry)
	cacheMu      sync.RWMutex
	cacheTTL     = 2 * time.Second
)

func getCached(key string) (any, bool) {
	cacheMu.RLock()
	defer cacheMu.RUnlock()
	entry, ok := rankingCache[key]
	if !ok || time.Now().After(entry.expiresAt) {
		return nil, false
	}
	return entry.data, true
}

func setCache(key string, data any) {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	rankingCache[key] = cacheEntry{data: data, expiresAt: time.Now().Add(cacheTTL)}
}

func InvalidateRankingCache() {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	rankingCache = make(map[string]cacheEntry)
}

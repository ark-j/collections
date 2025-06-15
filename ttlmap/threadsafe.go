package ttlmap

import (
	"sync"
	"time"
)

// TTLMap is thread safe time to live map
type TTLMap[K comparable, V any] struct {
	sync.RWMutex
	entries map[K]value[V]
	flush   time.Duration
	close   chan struct{}
}

func New[K comparable, V any](size int, flush time.Duration) *TTLMap[K, V] {
	tm := &TTLMap[K, V]{
		entries: make(map[K]value[V], size),
		flush:   flush,
		close:   make(chan struct{}),
	}
	go tm.flusher()
	return tm
}

// Set adds or override the key value pair along with it's time to live value.
func (tm *TTLMap[K, V]) Set(key K, val V, ttl time.Duration) {
	tm.Lock()
	tm.entries[key] = value[V]{
		val: val,
		ttl: time.Now().Add(ttl).UnixNano(),
	}
	tm.Unlock()
}

// Get returns the value for specified key. If value is expired or not in map
// the default value will be returned.
//
//nolint:ireturn
func (tm *TTLMap[K, V]) Get(key K) (V, bool) {
	tm.RLock()
	v, ok := tm.entries[key]
	tm.RUnlock()
	if !ok || time.Now().UnixNano() > v.ttl {
		var v V
		return v, false
	}
	return v.val, true
}

// Delete removes the key value pair for provided key
func (tm *TTLMap[K, V]) Delete(key K) {
	tm.Lock()
	delete(tm.entries, key)
	tm.Unlock()
}

// Len returns number of keys present which are not expired
func (tm *TTLMap[K, V]) Len() int {
	now := time.Now().UnixNano()
	count := 0
	tm.RLock()
	for _, v := range tm.entries {
		if now < v.ttl {
			count++
		}
	}
	tm.RUnlock()
	return count
}

// TotalLen returns total number of keys even expired keys are included.
func (tm *TTLMap[K, V]) TotalLen() int {
	tm.RLock()
	defer tm.RUnlock()
	return len(tm.entries)
}

// Close gracefully close the TTLMap.
// WARN: do not use map after closing as it may not work as intended.
func (tm *TTLMap[K, V]) Close() {
	tm.close <- struct{}{}
	close(tm.close)
}

// RemoveExpired will remove the expired keys from underlying map
func (tm *TTLMap[K, V]) RemoveExpired() {
	now := time.Now().UnixNano()
	tm.Lock()
	for k, v := range tm.entries {
		if now > v.ttl {
			delete(tm.entries, k)
		}
	}
	tm.Unlock()
}

// flusher periodically removes the expired keys.
func (tm *TTLMap[K, V]) flusher() {
	ticker := time.NewTicker(tm.flush)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if tm.TotalLen() > 0 {
				tm.RemoveExpired()
			}
		case <-tm.close:
			return
		}
	}
}

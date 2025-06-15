package ttlmap

import "time"

// value is a internal struct to represent value along with it's expirey
type value[V any] struct {
	val V
	ttl int64
}

// TTLMap is thread safe time to live map
type UnsafeTTLMap[K comparable, V any] struct {
	entries map[K]value[V]
	flush   time.Duration
	close   chan struct{}
}

func NewUnsafe[K comparable, V any](size int, flush time.Duration) *UnsafeTTLMap[K, V] {
	tm := &UnsafeTTLMap[K, V]{
		entries: make(map[K]value[V], size),
		flush:   flush,
		close:   make(chan struct{}),
	}
	go tm.flusher()
	return tm
}

// Set adds or override the key value pair along with it's time to live value.
func (tm *UnsafeTTLMap[K, V]) Set(key K, val V, ttl time.Duration) {
	tm.entries[key] = value[V]{
		val: val,
		ttl: time.Now().Add(ttl).UnixNano(),
	}
}

// Get returns the value for specified key. If value is
// expired or not in map the default value will be returned.
//
//nolint:ireturn
func (tm *UnsafeTTLMap[K, V]) Get(key K) (V, bool) {
	v, ok := tm.entries[key]
	if !ok || time.Now().UnixNano() > v.ttl {
		var v V
		return v, false
	}
	return v.val, true
}

// Delete removes the key value pair for provided key
func (tm *UnsafeTTLMap[K, V]) Delete(key K) {
	delete(tm.entries, key)
}

// Len returns number of keys present which are not expired
func (tm *UnsafeTTLMap[K, V]) Len() int {
	now := time.Now().UnixNano()
	count := 0
	for _, v := range tm.entries {
		if now < v.ttl {
			count++
		}
	}
	return count
}

// TotalLen returns total number of keys even expired keys are included.
func (tm *UnsafeTTLMap[K, V]) TotalLen() int {
	return len(tm.entries)
}

// Close gracefully close the TTLMap.
// WARN: do not use map after closing as it may not work as intended.
func (tm *UnsafeTTLMap[K, V]) Close() {
	tm.close <- struct{}{}
	close(tm.close)
}

// RemoveExpired will remove the expired keys from underlying map
func (tm *UnsafeTTLMap[K, V]) RemoveExpired() {
	now := time.Now().UnixNano()
	for k, v := range tm.entries {
		if now > v.ttl {
			delete(tm.entries, k)
		}
	}
}

// flusher periodically removes the expired keys.
func (tm *UnsafeTTLMap[K, V]) flusher() {
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

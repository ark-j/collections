// package smap provide courrent safe generic map
package smap

import (
	"maps"
	"sync"
)

type Map[K comparable, V any] struct {
	m  map[K]V
	mu sync.RWMutex
}

func New[K comparable, V any](capacity int) *Map[K, V] {
	return &Map[K, V]{
		m: make(map[K]V, capacity),
	}
}

// Put adds key-val entry to map
func (m *Map[K, V]) Put(key K, val V) {
	m.mu.Lock()
	m.m[key] = val
	m.mu.Unlock()
}

// Get get value provided key
func (m *Map[K, V]) Get(key K) V { //nolint
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.m[key]
}

// Keys returns all the ket present in map
func (m *Map[K, V]) Keys() []K {
	keys := make([]K, 0, m.Len())
	m.mu.RLock()
	for k := range m.m {
		keys = append(keys, k)
	}
	m.mu.RUnlock()
	return keys
}

// Vals returns all the values present in map
func (m *Map[K, V]) Vals() []V {
	vals := make([]V, 0, m.Len())
	m.mu.RLock()
	for _, v := range m.m {
		vals = append(vals, v)
	}
	m.mu.RUnlock()
	return vals
}

// Len return length of map
func (m *Map[K, V]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.m)
}

// Delete deletes the pair based in key
func (m *Map[K, V]) Delete(key K) {
	m.mu.Lock()
	delete(m.m, key)
	m.mu.Unlock()
}

// DeleteFunc is mimic [maps.DeleteFunc]
// you only need to provide func as argumnet
func (m *Map[K, V]) DeleteFunc(fn func(k K, v V) bool) {
	m.mu.Lock()
	maps.DeleteFunc(m.m, fn)
	m.mu.Unlock()
}

// Clear deletes all the entries in map
func (m *Map[K, V]) Clear() {
	m.mu.Lock()
	clear(m.m)
	m.mu.Unlock()
}

// Clone returns the clone of maps
// beware for passing pointers as values
// because they can still modify original map
func (m *Map[K, V]) Clone() map[K]V {
	cm := make(map[K]V, m.Len())
	m.mu.RLock()
	for k, v := range m.m {
		cm[k] = v
	}
	m.mu.RUnlock()
	return cm
}

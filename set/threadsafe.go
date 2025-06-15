package set

import (
	"encoding/json"
	"sync"
)

type Set[T comparable] struct {
	mu sync.RWMutex
	us UnsafeSet[T]
}

func NewSafe[T comparable](size int) *Set[T] {
	return &Set[T]{us: NewUnsafe[T](size)}
}

func NewSafeFromSlice[T comparable](slice []T) *Set[T] {
	return &Set[T]{}
}

func (s *Set[T]) Add(v T) {
	s.mu.Lock()
	s.us.Add(v)
	s.mu.Unlock()
}

func (s *Set[T]) Contains(v T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.us.Contains(v)
}

func (s *Set[T]) Append(v ...T) {
	s.mu.Lock()
	s.us.Append(v...)
	s.mu.Unlock()
}

func (s *Set[T]) Cardinality() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.us)
}

func (s *Set[T]) Clear() {
	s.mu.Lock()
	s.us.Clear()
	s.mu.Unlock()
}

func (s *Set[T]) Clone() *Set[T] {
	s.mu.RLock()
	us := s.us.Clone()
	s.mu.RUnlock()
	return &Set[T]{us: us}
}

func (s *Set[T]) Equal(o *Set[T]) bool {
	s.mu.RLock()
	o.mu.RLock()
	defer func() {
		s.mu.RUnlock()
		o.mu.RUnlock()
	}()
	return s.us.Equal(o.us)
}

func (s *Set[T]) Intersect(o *Set[T]) *Set[T] {
	s.mu.RLock()
	o.mu.RLock()
	us := s.us.Intersect(o.us)
	s.mu.RUnlock()
	o.mu.RUnlock()
	return &Set[T]{us: us}
}

func (s *Set[T]) Union(o *Set[T]) *Set[T] {
	s.mu.RLock()
	o.mu.RLock()
	us := s.us.Union(o.us)
	s.mu.RUnlock()
	o.mu.RUnlock()
	return &Set[T]{us: us}
}

func (s *Set[T]) Difference(o *Set[T]) *Set[T] {
	s.mu.RLock()
	o.mu.RLock()
	us := s.us.Difference(o.us)
	s.mu.RUnlock()
	o.mu.RUnlock()
	return &Set[T]{us: us}
}

func (s *Set[T]) IsEmpty() bool {
	return s.Cardinality() == 0
}

func (s *Set[T]) IsSubset(o *Set[T]) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	o.mu.RLock()
	defer o.mu.RUnlock()
	return s.us.IsSubset(o.us)
}

func (s *Set[T]) IsSuperset(other *Set[T]) bool {
	return other.IsSubset(s)
}

func (s *Set[T]) Remove(v T) {
	s.mu.Lock()
	delete(s.us, v)
	s.mu.Unlock()
}

func (s *Set[T]) RemoveMulti(v ...T) {
	s.mu.Lock()
	s.us.RemoveMulti(v...)
	s.mu.Unlock()
}

func (s *Set[T]) FromSlice(slice []T) {
	for _, k := range slice {
		s.us.Add(k)
	}
}

func (s *Set[T]) ToSlice() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.us.ToSlice()
}

func (s *Set[T]) MarshalJSON() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return json.Marshal(s.ToSlice())
}

func (s *Set[T]) UnmarshalJSON(data []byte) error {
	var arr []T
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	s.FromSlice(arr)
	return nil
}

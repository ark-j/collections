package set

import "encoding/json"

type UnsafeSet[T comparable] map[T]struct{}

func NewUnsafe[T comparable](size int) UnsafeSet[T] {
	return make(UnsafeSet[T], size)
}

func NewUnsafeFromSlice[T comparable](slice []T) UnsafeSet[T] {
	s := make(UnsafeSet[T], len(slice))
	s.FromSlice(slice)
	return s
}

func (s UnsafeSet[T]) Add(v T) {
	s[v] = struct{}{}
}

func (s UnsafeSet[T]) Contains(v T) bool {
	_, ok := s[v]
	return ok
}

func (s UnsafeSet[T]) Append(v ...T) {
	for _, val := range v {
		s[val] = struct{}{}
	}
}

func (s UnsafeSet[T]) Cardinality() int {
	return len(s)
}

func (s UnsafeSet[T]) Clear() {
	for k := range s {
		delete(s, k)
	}
}

func (s UnsafeSet[T]) Clone() UnsafeSet[T] {
	cloned := NewUnsafe[T](s.Cardinality())
	for v := range s {
		cloned.Add(v)
	}
	return cloned
}

func (s UnsafeSet[T]) Equal(other UnsafeSet[T]) bool {
	if s.Cardinality() != other.Cardinality() {
		return false
	}
	for v := range s {
		if !other.Contains(v) {
			return false
		}
	}
	return true
}

func (s UnsafeSet[T]) Intersect(other UnsafeSet[T]) UnsafeSet[T] {
	intersect := NewUnsafe[T](0)
	if s.Cardinality() < other.Cardinality() {
		for v := range s {
			if other.Contains(v) {
				intersect.Add(v)
			}
		}
	}
	return intersect
}

func (s UnsafeSet[T]) Union(other UnsafeSet[T]) UnsafeSet[T] {
	union := NewUnsafe[T](0)
	union.FromSlice(s.ToSlice())
	union.FromSlice(other.ToSlice())
	return union
}

func (s UnsafeSet[T]) Difference(other UnsafeSet[T]) UnsafeSet[T] {
	diff := NewUnsafe[T](0)
	for v := range s {
		if !other.Contains(v) {
			diff.Add(v)
		}
	}
	return diff
}

func (s UnsafeSet[T]) IsEmpty() bool {
	return s.Cardinality() == 0
}

func (s UnsafeSet[T]) IsSubset(other UnsafeSet[T]) bool {
	if s.Cardinality() > other.Cardinality() {
		return false
	}
	for v := range s {
		if !s.Contains(v) {
			return false
		}
	}
	return true
}

func (s UnsafeSet[T]) IsSuperset(other UnsafeSet[T]) bool {
	return other.IsSubset(s)
}

func (s UnsafeSet[T]) Remove(v T) {
	delete(s, v)
}

func (s UnsafeSet[T]) RemoveMulti(v ...T) {
	for _, val := range v {
		delete(s, val)
	}
}

func (s UnsafeSet[T]) FromSlice(slice []T) {
	for _, k := range slice {
		s.Add(k)
	}
}

func (s UnsafeSet[T]) ToSlice() []T {
	slice := make([]T, 0, s.Cardinality())
	for v := range s {
		slice = append(slice, v)
	}
	return slice
}

func (s UnsafeSet[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.ToSlice())
}

func (s UnsafeSet[T]) UnmarshalJSON(data []byte) error {
	var arr []T
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	s.FromSlice(arr)
	return nil
}

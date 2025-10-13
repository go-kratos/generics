package generics

import "sync"

// List is a thread-safe generic slice-based list.
// It uses RWMutex to ensure safe concurrent reads and writes.
type List[T any] struct {
	mu   sync.RWMutex
	data []T
}

// NewList creates a new List with optional initial elements.
func NewList[T any](items ...T) *List[T] {
	d := make([]T, 0, len(items))
	d = append(d, items...)
	return &List[T]{data: d}
}

// Len returns the number of elements in the list.
func (l *List[T]) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.data)
}

// Clear removes all elements from the list.
func (l *List[T]) Clear() {
	l.mu.Lock()
	l.data = nil
	l.mu.Unlock()
}

// Append adds items to the end of the list.
func (l *List[T]) Append(items ...T) *List[T] {
	if len(items) == 0 {
		return l
	}
	l.mu.Lock()
	l.data = append(l.data, items...)
	l.mu.Unlock()
	return l
}

// Get returns the item at index i.
// It returns false if i is out of bounds.
func (l *List[T]) Get(i int) (T, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if i < 0 || i >= len(l.data) {
		var zero T
		return zero, false
	}
	return l.data[i], true
}

// Set replaces the item at index i with value.
// It returns false if i is out of bounds.
func (l *List[T]) Set(i int, value T) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if i < 0 || i >= len(l.data) {
		return false
	}
	l.data[i] = value
	return true
}

// RemoveAt removes and returns the item at index i.
// It returns false if i is out of bounds.
func (l *List[T]) RemoveAt(i int) (T, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if i < 0 || i >= len(l.data) {
		var zero T
		return zero, false
	}
	v := l.data[i]
	l.data = append(l.data[:i], l.data[i+1:]...)
	return v, true
}

// Range iterates over a snapshot of the list.
// The callback receives the index and item. If it returns false, iteration stops.
func (l *List[T]) Range(f func(index int, item T) bool) {
	l.mu.RLock()
	cpy := make([]T, len(l.data))
	copy(cpy, l.data)
	l.mu.RUnlock()
	for i, v := range cpy {
		if !f(i, v) {
			break
		}
	}
}

// ToSlice returns a copy of the underlying slice.
// Safe for concurrent use.
func (l *List[T]) ToSlice() []T {
	l.mu.RLock()
	defer l.mu.RUnlock()
	cpy := make([]T, len(l.data))
	copy(cpy, l.data)
	return cpy
}

// Clone creates and returns a shallow copy of the list.
func (l *List[T]) Clone() *List[T] {
	l.mu.RLock()
	defer l.mu.RUnlock()
	c := make([]T, len(l.data))
	copy(c, l.data)
	return &List[T]{data: c}
}

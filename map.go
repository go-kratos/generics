package generics

import "sync"

// Map is a concurrent map with generic key and value types.
type Map[K comparable, V any] struct {
	sync.Map
}

// NewMap creates and returns a new Map instance.
func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{}
}

// Clear removes all entries from the map.
func (m *Map[K, V]) Clear() {
	m.Map.Clear()
}

// CompareAndDelete deletes the entry for a key only if it is currently mapped to a given value.
func (m *Map[K, V]) CompareAndDelete(key K, value V) (deleted bool) {
	return m.Map.CompareAndDelete(key, value)
}

// CompareAndSwap swaps the entry for a key only if it is currently mapped to a given value.
func (m *Map[K, V]) CompareAndSwap(key, old, new any) (swapped bool) {
	return m.Map.CompareAndSwap(key, old, new)
}

// Delete removes the value for a given key.
func (m *Map[K, V]) Delete(key K) {
	m.Map.Delete(key)
}

// Load retrieves the value for a given key.
func (m *Map[K, V]) Load(key K) (V, bool) {
	value, ok := m.Map.Load(key)
	if !ok {
		var zero V
		return zero, false
	}
	return value.(V), true
}

// LoadAndDelete retrieves and deletes the value for a given key.
func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	loadedValue, ok := m.Map.LoadAndDelete(key)
	if !ok {
		var zero V
		return zero, false
	}
	return loadedValue.(V), true
}

// LoadOrStore retrieves the existing value for a key or stores and returns the given value if the key is not present.
func (m *Map[K, V]) LoadOrStore(key K, value V) (V, bool) {
	loaded, ok := m.Map.LoadOrStore(key, value)
	if !ok {
		return value, false
	}
	return loaded.(V), true
}

// Range iterates over all key-value pairs in the map.
func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	m.Map.Range(func(key, value any) bool {
		return f(key.(K), value.(V))
	})
}

// Store sets the value for a given key.
func (m *Map[K, V]) Store(key K, value V) {
	m.Map.Store(key, value)
}

// Swap sets the value for a key and returns the previous value and whether it was present.
func (m *Map[K, V]) Swap(key, value any) (previous any, loaded bool) {
	return m.Map.Swap(key, value)
}

// Clone creates and returns a shallow copy of the map as a standard map.
func (m *Map[K, V]) ToMap() map[K]V {
	clone := make(map[K]V)
	m.Range(func(key K, value V) bool {
		clone[key] = value
		return true
	})
	return clone
}

// Clone creates and returns a shallow copy of the Map.
func (m *Map[K, V]) Clone() *Map[K, V] {
	clone := NewMap[K, V]()
	m.Range(func(key K, value V) bool {
		clone.Store(key, value)
		return true
	})
	return clone
}

  # Generic Map

  A tiny, type-safe wrapper around Go’s sync.Map using generics. It provides a typed API for concurrent maps so you can avoid manual casting while keeping sync.Map’s performance characteristics and
  semantics.

  - Type-safe keys and values via generics
  - Thin wrapper over sync.Map for familiar behavior
  - Convenience methods like Clone and Clear

  ## Installation

  - go get github.com/go-kratos/syncmap
  - Import: import "github.com/go-kratos/syncmap"

  Requires a Go version with generics support (Go 1.18+). The module currently declares go 1.24 in go.mod.

  ## Quick Start
  package main

  import (
  	"fmt"

  	"github.com/go-kratos/syncmap"
  )

  func main() {
  	var m syncmap.Map[string, int]

  	// Store and Load
  	m.Store("a", 1)
  	if v, ok := m.Load("a"); ok {
  		fmt.Println("a =", v) // a = 1
  	}

  	// LoadOrStore
  	v, loaded := m.LoadOrStore("b", 42)
  	fmt.Println("b =", v, "loaded =", loaded) // b = 42 loaded = false

  	// CompareAndSwap (parity with sync.Map)
  	swapped := m.CompareAndSwap("b", 42, 43)
  	fmt.Println("b swapped =", swapped) // true

  	// Swap (returns previous value, if any)
  	prev, existed := m.Swap("a", 2)
  	fmt.Println("prev a =", prev, "existed =", existed) // prev a = 1 existed = true

  	// Range over entries
  	m.Range(func(k string, v int) bool {
  		fmt.Printf("%s => %d\n", k, v)
  		return true
  	})

  	// Clone to a regular map snapshot
  	clone := m.Clone()
  	fmt.Println("clone[a] =", clone["a"]) // clone[a] = 2

  	// Delete and Clear
  	m.Delete("a")
  	m.Clear()
  }
  ## API

  Typed wrapper methods mirroring sync.Map:

  - func (m *Map[K, V]) Store(key K, value V)
  - func (m *Map[K, V]) Load(key K) (V, bool)
  - func (m *Map[K, V]) LoadOrStore(key K, value V) (V, bool)
  - func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool)
  - func (m *Map[K, V]) Delete(key K)
  - func (m *Map[K, V]) Range(f func(key K, value V) bool)
  - func (m *Map[K, V]) Clear()
  - func (m *Map[K, V]) CompareAndDelete(key K, value V) (deleted bool)
  - func (m *Map[K, V]) CompareAndSwap(key, old, new any) (swapped bool)  // parity with sync.Map
  - func (m *Map[K, V]) Swap(key, value any) (previous any, loaded bool)  // parity with sync.Map
  - func (m *Map[K, V]) Clone() map[K]V  // shallow snapshot

  Notes:

  - CompareAndSwap and Swap use any parameters to match sync.Map signatures.
  - Clone returns a snapshot; it doesn’t stay in sync with the original map.

  ## Usage Notes

  - Semantics are identical to sync.Map: best suited for highly concurrent, read-heavy workloads and dynamically changing keys.
  - For simple, write-heavy cases with stable keys, a regular map with sync.RWMutex may be more efficient.

  ## License

  MIT License. See LICENSE.

  Would you like me to create README.md with this content in the repository? If yes, I’ll add the file now.

# go-kratos/generics

Lightweight, type-safe generic containers for Go with concurrency-friendly APIs.

Included containers:

- Set[T] — a generic set built on top of Map.
- Map[K,V] — a type-safe wrapper around `sync.Map`.
- Slice[T] — a thread-safe, slice-backed list (uses `sync.RWMutex`).

Great for simple, robust concurrent read/write scenarios. Zero dependencies, easy to integrate.

## Features

- Type-safe APIs via Go generics.
- Concurrency aware: Slice uses `RWMutex`; Map wraps `sync.Map`.
- Snapshot operations: `Slice.Range` iterates over a snapshot; `ToSlice`/`ToMap` return copies.
- Familiar, minimal APIs: `NewSlice`/`NewMap`/`NewSet`, plus `Clone`, `Clear`, etc.
- Standard library only, no external deps.

Requirements: Go 1.19+ (generics + newer `sync.Map` APIs).

## Install

```bash
go get github.com/go-kratos/generics@latest
```

## Quick Start

### Set[T]

```go
package main

import (
    "fmt"
    "github.com/go-kratos/generics"
)

func main() {
    s := generics.NewSet[string]("a", "b")
    s.Insert("c").Delete("b")

    fmt.Println(s.Has("a"))         // true
    fmt.Println(s.HasAny("x", "c")) // true
    fmt.Println(s.HasAll("a", "c")) // true

    t := s.Clone()
    fmt.Println(t.HasAll("a", "c")) // true
}
```

Common methods: `Insert`, `Delete`, `Has`, `HasAny`, `HasAll`, `Clear`, `Clone`.

### Map[K,V]

```go
package main

import (
    "fmt"
    "github.com/go-kratos/generics"
)

func main() {
    m := generics.NewMap[string, int]()
    m.Store("a", 1)

    if v, ok := m.Load("a"); ok {
        fmt.Println(v) // 1
    }

    v, ok = m.LoadOrStore("b", 2)

    m.Range(func(k string, v int) bool {
        fmt.Println(k, v)
        return true
    })

    // copy into a built-in map
    fmt.Println(m.ToMap())
}
```

Common methods: `Store`, `Load`, `LoadOrStore`, `LoadAndDelete`, `Delete`, `Clear`, `Range`, `ToMap`, `Clone`.

### Slice[T]

```go
package main

import (
    "fmt"
    "github.com/go-kratos/generics"
)

func main() {
    l := generics.NewSlice[int](1, 2, 3)
    l.Append(4, 5)  // append 4, 5
    l.Insert(1, 99) // insert at index

    if v, ok := l.Get(1); ok {
        fmt.Println(v) // 99
    }

    // iterate over a snapshot
    l.Range(func(i int, v int) bool {
        fmt.Printf("%d:%d ", i, v)
        return true
    })
    fmt.Println()

    fmt.Println(l.ToSlice())
}
```

Common methods: `Append`, `Get`, `Set`, `RemoveAt`, `Range`, `Slice/SliceStart/SliceEnd`, `ToSlice`, `Clone`, `Len`, `Clear`.



## Concurrency Notes

- Slice: write ops are mutex-protected; `Range`/`ToSlice` work on a snapshot to avoid long-held locks.
- Map: type-safe wrapper over `sync.Map`; good for read-heavy or cross-goroutine sharing.
- Set: built on top of the concurrent Map; methods are safe for concurrent use.

Note: `Slice.Range` copies the underlying slice; for very large lists, consider memory impact. `ToSlice`/`ToMap` similarly return copies.

## Design Choices

- Prefer simplicity: keep the surface area small and familiar.
- Readability first: names are aligned with standard library conventions.
- Reliability: built solely on standard library concurrency primitives.

## Roadmap

- Additional containers: queue, stack, ring buffer, etc.
- More helpers and converters.
- More docs and examples.

Contributions via Issues/PRs are welcome!

## License

MIT License. See `LICENSE` for details.

# go-kratos/kit

轻量、类型安全的通用工具包，包含并发友好的泛型容器与简洁的重试器。

内置包：

- containers/maps：基于 `sync.Map` 的类型安全泛型 Map。
- containers/sets：基于 Map 的泛型 Set。
- containers/slices：使用 `sync.RWMutex` 封装的并发安全 Slice 列表。
- retry：带指数退避的通用重试器，可配置重试条件与退避参数。

仅依赖标准库，易于集成到任意项目。

## 安装

```bash
go get github.com/go-kratos/kit@latest
```

按需引入子包，例如：

```go
import (
    "github.com/go-kratos/kit/containers/maps"
    "github.com/go-kratos/kit/containers/sets"
    "github.com/go-kratos/kit/containers/slices"
    "github.com/go-kratos/kit/retry"
)
```

## 快速开始

### containers/sets.Set[T]

```go
package main

import (
    "fmt"
    "github.com/go-kratos/kit/containers/sets"
)

func main() {
    s := sets.New[string]("a", "b")
    s.Insert("c").Delete("b")

    fmt.Println(s.Has("a"))          // true
    fmt.Println(s.HasAny("x", "c"))  // true
    fmt.Println(s.HasAll("a", "c"))  // true

    t := s.Clone()
    fmt.Println(t.HasAll("a", "c"))  // true

    fmt.Println(s.ToSlice())           // [a c]（顺序不保证）
}
```

常用方法：`Insert`、`Delete`、`Has`、`HasAny`、`HasAll`、`Clear`、`Clone`、`ToSlice`。

JSON 支持：Set 会被编码为元素数组；解码时会填充集合。

### containers/maps.Map[K,V]

```go
package main

import (
    "fmt"
    "github.com/go-kratos/kit/containers/maps"
)

func main() {
    m := maps.New[string, int]()
    m.Store("a", 1)

    if v, ok := m.Load("a"); ok {
        fmt.Println(v) // 1
    }

    if v, loaded := m.LoadOrStore("b", 2); !loaded {
        fmt.Println(v) // 2（首次存入）
    }

    m.Range(func(k string, v int) bool {
        fmt.Println(k, v)
        return true
    })

    // 拷贝为内建 map
    fmt.Println(m.ToMap())

    // 克隆为新的并发 Map
    mm := m.Clone()
    _ = mm
}
```

常用方法：`Store`、`Load`、`LoadOrStore`、`LoadAndDelete`、`Delete`、`Clear`、`Range`、`ToMap`、`Clone`。

JSON 支持：直接序列化/反序列化为对象（map）。

### containers/slices.Slice[T]

```go
package main

import (
    "fmt"
    "github.com/go-kratos/kit/containers/slices"
)

func main() {
    l := slices.New[int](1, 2, 3)
    l.Append(4, 5)

    if v, ok := l.Get(1); ok {
        fmt.Println(v) // 2
    }

    // 设置与删除
    _ = l.Set(1, 99)
    if v, ok := l.RemoveAt(0); ok {
        fmt.Println("removed:", v) // 1
    }

    // 快照遍历
    l.Range(func(i int, v int) bool {
        fmt.Printf("%d:%d ", i, v)
        return true
    })
    fmt.Println()

    fmt.Println(l.ToSlice()) // 返回副本
}
```

常用方法：`Append`、`Get`、`Set`、`RemoveAt`、`Range`、`Slice`/`SliceStart`/`SliceEnd`、`ToSlice`、`Clone`、`Len`、`Clear`。

JSON 支持：序列化/反序列化为数组；内部做并发保护。

### retry 重试器

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "time"
    "github.com/go-kratos/kit/retry"
)

func main() {
    // 最多重试 3 次，且仅对临时错误重试
    r := retry.New(
        3,
        retry.WithRetryable(func(err error) bool { return errors.Is(err, context.DeadlineExceeded) }),
        retry.WithBaseDelay(10*time.Millisecond),
        retry.WithMaxDelay(1*time.Second),
        retry.WithMultiplier(1.6),
        retry.WithJitter(0.2),
    )

    err := r.Do(context.Background(), func(ctx context.Context) error {
        // 你的业务逻辑
        return context.DeadlineExceeded
    })

    fmt.Println("done:", err)
}
```

便捷方法：

- `retry.Do(ctx, fn)`：使用默认配置重试。
- `retry.Infinite(ctx, fn)`：无限次重试（直到成功或 `ctx` 结束）。

可调参数：`WithBaseDelay`、`WithMaxDelay`、`WithMultiplier`、`WithJitter`、`WithRetryable`。

## 并发与性能说明

- slices：写操作使用互斥锁保护；`Range`、`ToSlice` 基于快照，避免长时间持锁。
- maps：类型安全封装 `sync.Map`，适合读多写少或跨 goroutine 共享场景。
- sets：基于并发 Map 构建，方法并发安全。

注意：`Slice.Range` 会复制底层切片；超大列表遍历时需考虑内存开销。`ToSlice`/`ToMap` 均返回副本。

## 版本要求

需要支持泛型的 Go 版本（Go 1.18+）。

## 许可证

MIT，详见 `LICENSE`。

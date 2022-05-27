# Golang 基础知识与框架体系 part 4: 等待组 waitGroup


这篇文章总结了 Golang 的知识体系: waitGroup，包括其中的底层实现等等。

<!--more-->

## 等待组 WaitGroup

很多情况下，我们正需要知道 goroutine 是否完成。这需要借助 `sync` 包的 `WaitGroup` 来实现。`WaitGroup` 是 `sync` 包中的一个 `struct` 类型，用来收集需要等待执行完成的 goroutine。下面是它的定义：

```go
type WaitGroup struct {
        // Has unexported fields.
}
  // A WaitGroup waits for a collection of goroutines to finish. The main
  // goroutine calls Add to set the number of goroutines to wait for. Then each
  // of the goroutines runs and calls Done when finished. At the same time, Wait
  // can be used to block until all goroutines have finished.

  // A WaitGroup must not be copied after first use.


func (wg *WaitGroup) Add(delta int)
func (wg *WaitGroup) Done()
func (wg *WaitGroup) Wait()
```

它有3个方法：

- `Add()`：每次激活想要被等待完成的 goroutine 之前，先调用 `Add()`，用来设置或添加要等待完成的 goroutine 数量例如 `Add(2)` 或者两次调用 `Add(1)` 都会设置等待计数器的值为2，表示要等待2个 goroutine 完成
- `Done()`：每次需要等待的 goroutine 在真正完成之前，应该调用该方法来人为表示 goroutine 完成了，该方法会对等待计数器减1
- `Wait()`：在等待计数器减为0之前，`Wait()` 会一直阻塞当前的 goroutine

```go
package main

import (  
    "fmt"
    "sync"
    "time"
)

func process(i int, wg *sync.WaitGroup) {  
    fmt.Println("started Goroutine ", i)
    time.Sleep(2 * time.Second)
    fmt.Printf("Goroutine %d ended\n", i)
    wg.Done()
}

func main() {  
    no := 3
    var wg sync.WaitGroup
    for i := 0; i < no; i++ {
        wg.Add(1)
        go process(i, &wg)
    }
    wg.Wait()
    fmt.Println("All go routines finished executing")
}
```
上面激活了3个 goroutine，每次激活 goroutine 之前，都先调用 `Add()` 方法增加一个需要等待的 goroutine 计数。每个 goroutine 都运行 `process()` 函数，这个函数在执行完成时需要调用 `Done()` 方法来表示 goroutine 的结束。激活3个 goroutine 后，main goroutine 会执行到 `Wait()`，由于每个激活的 goroutine 运行的 `process()` 都需要睡眠2秒，所以 main goroutine 在 `Wait()` 这里会阻塞一段时间(大约2秒)，当所有 goroutine 都完成后，等待计数器减为0，`Wait()` 将不再阻塞，于是 `main` goroutine 得以执行后面的 `Println()`。

{{< admonition >}}
还有一点需要特别注意的是 `process()` 中使用指针类型的 `*sync.WaitGroup` 作为参数，这里不能使用值类型的 `sync.WaitGroup` 作为参数，因为这意味着每个 goroutine 都拷贝一份 `wg`，每个 goroutine 都使用自己的 `wg`。这显然是不合理的，这3个 goroutine 应该共享一个 `wg`，才能知道这3个 goroutine 都完成了。实际上，如果使用值类型的参数，main goroutine 将会永久阻塞而导致产生死锁。
{{< /admonition >}}


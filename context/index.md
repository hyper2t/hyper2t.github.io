# Go 基础知识与框架体系 part 5: context包


这篇文章总结了 Go 的知识体系:context包，包括其中的底层实现等等。

<!--more-->

## context 包

`context` 包主要用于父子任务之间的同步取消信号，本质上是一种协程调度的方式。可以通过简化对于处理单个请求的多个 Goroutine 之间与请求域的数据、超时和退出等操作。来看两个例子：

```go
package main
import(
	"context"
	"sync"
	"fmt"
	"time"
)
// 我们定义一个worker function
func worker(ctx context.Context, wg *sync.WaitGroup) error {
  defer wg.Done()

  for {
    select {
    case <- ctx.Done():
      return ctx.Err()
    default:
      fmt.Println("hello")
    }
  }
}

func main(){
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

  var wg sync.WaitGroup
  for i := 0; i < 10; i++{
    wg.Add(1)
    go worker(ctx, &wg)
  }
  time.Sleep(time.Second)
  cancel()

  wg.Wait()
}

```

当并发体超时或者 `main` 主动停止 `worker` goroutine 时，`worker` goroutine 都可以安全退出。

另外在使用 `context` 时有两点值得注意：上游任务仅仅使用 `context` 通知下游任务不再需要，但不会直接干涉和中断下游任务的执行，由下游任务自行决定后续的处理操作，也就是说 `context` 的取消操作是无侵入的；`context` 是线程安全的，因为 `context` 本身是不可变的（immutable），因此可以放心地在多个协程中传递使用。


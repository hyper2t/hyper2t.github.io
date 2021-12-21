# Golang 基础知识与框架体系


这篇文章总结了 Golang 的知识体系和实际应用中用到的框架，包括其中的底层实现等等。

<!--more-->

{{< admonition >}}
说起 Golang， 大家都会第一时间想到高并发和 Golang 作为主流的后端开发语言的优势，本文主要讲 Golang 主要知识体系，包括**数组和切片**、**协程的调度**原理、等待组 **waitGroup**、**channel** 的底层实现、互斥锁 **mutex** 的实现、**interface** 中的多态等等。
{{< /admonition >}}

## 1 数组和切片

### 1. 切片的本质

切片的本质就是对底层数组的封装，它包含了三个信息

- 底层数组的指针
- 切片的长度(len)
- 切片的容量(cap)

切片的容量指的是**数组中的头指针指向的位置至数组最后一位的长度**。举个例子，现在有一个数组 `a := [8]int {0,1,2,3,4,5,6,7}`，切片 `s1 := a[:5]`，相应示意图如下

![切片 s1 和底层数组 a](/slice1.png "图1: 切片 s1 和底层数组 a")

切片 `s2 := a[3:6]`，相应示意图如下：

![切片 s2 和底层数组 a](/slice2.png "图2：切片 s2 和底层数组 a")

### 2. 切片的扩容

Go 中切片扩容的策略是这样的：

- 如果切片的容量小于 1024 个元素，于是扩容的时候就翻倍增加容量。上面那个例子也验证了这一情况，总容量从原来的4个翻倍到现在的8个。
- 一旦元素个数超过 1024 个元素，那么增长因子就变成 1.25 ，即每次增加原来容量的四分之一。

{{< admonition type=tip open=true >}}
注意：扩容扩大的容量都是针对原来的容量而言的，而不是针对原来数组的长度而言的。
{{< /admonition >}}

举例

{{< admonition type=example >}}
```go
// make()函数创建切片
fmt.Println()
var slices = make([]int, 4, 8)
//[0 0 0 0]
fmt.Println(slices)
// 长度：4, 容量8
fmt.Printf("长度：%d, 容量%d", len(slices), cap(slices))
```
{{< /admonition >}}

需要注意的是，golang 中没办法通过下标来给切片扩容，如果需要扩容，需要用到 `append`

```go
slices2 := []int{1,2,3,4}
slices2 = append(slices2, 5)
fmt.Println(slices2)
// 输出结果 [1 2 3 4 5]
```

同时切片还可以将两个切片进行合并

```go
// 合并切片
slices3 := []int{6,7,8}
slices2 = append(slices2, slices3...)
fmt.Println(slices2)
// 输出结果  [1 2 3 4 5 6 7 8]
```



### 3. 切片的传递问题

切片本身传递给函数形参时是引用传递，但 `append` 后，切片长度变化时会被重新分配内存，而原来的切片还是指向原来地址，致使与初始状况传进来的地址不一样，要想对其值有改变操作，需使用指针类型操作。

我们来看一道 leetcode 78:

```go

package main

import "fmt"

func helper(nums []int, res *[][]int, tmp []int, level int) {
	if len(tmp) <= len(nums) {
		//长度一样的tmp用的是同一个地址，故可能会存在覆盖值得情况，
		// 长度不一样时重新开辟空间，将已有得元素复制进去
		//*res = append(*res, tmp) 如此处，最终长度为1的tmp会被最后3这个元素全部覆盖
		//以下相当于每次重新申请内存，使其指向的地址不一样，解决了最后地址一样的元素值被覆盖的状态状态
		var a []int
		a = append(a, tmp[:] ...)
		//res = append(res, a) 如果此处不是指针引用传递，在append后，res重新分配内存，与之前传进来的res地址不一样，最终res仍为空值
		*res = append(*res, a)
	}
	//fmt.Println(*res, "--->", tmp)
	for i := level; i < len(nums); i ++ {
		tmp = append(tmp, nums[i])
		helper(nums, res, tmp, i + 1)
		tmp = tmp[:len(tmp) - 1] //相当于删除tmp末位的最后一个元素
	}
}

func subsets(nums []int) [][]int {
	if len(nums) == 0 {
		return nil
	}
	var res [][]int
	var tmp []int
	helper(nums, &res, tmp, 0)
	return res
}


func main()  {
	pre := []int{1, 2, 3}
	fmt.Println(subsets(pre))
}
//错误结果：[[] [3] [1 3] [1 2 3] [1 3] [3] [2 3] [3]]， 可以看出，长度为1的切片都被3覆盖了，这由于它们的地址不一样
//正确输出：[[] [1] [1 2] [1 2 3] [1 3] [2] [2 3] [3]]， 这是因为每次都为a分配内存，其地址都与之前的不一样，故最终的值没有被覆盖
```
## 2 协程(coroutine)

在讲 goroutine 之前让我们先了解一下协程 (coroutine)，和线程类似，共享堆，不共享栈，协程的切换一般由程序员在代码中显式控制。它避免了上下文切换的额外耗费，兼顾了多线程的优点，简化了高并发程序的复杂。

### 1. golang 的协程(goroutine)

Goroutine 和其他语言的协程（coroutine）在使用方式上类似，但在区别上看，协程不是并发的，而Goroutine支持并发的。因此Goroutine可以理解为一种Go语言的协程。同时它可以运行在一个或多个线程上。来看个例子:

{{< admonition type=example >}}
```go
func Hello()  {
 fmt.Println("hello everybody , I'm Junqi Liao")
}

func main()  {
 go Hello()
 fmt.Println("Golang梦工厂")
}
```
{{< /admonition >}}

上面的程序，我们使用go又开启了一个 goroutine 执行 Hello 方法，但是我们运行这个程序，运行结果如下:
```
Golang梦工厂
```
这里出现这个问题的原因是我们启动的 goroutine 在 main 执行完就退出了。解决办法可以用 channel 通信让 goroutine 告诉 main 我执行完了，您再打印 “Golang梦工厂”。

```go
func Hello(ch chan int)  {
 fmt.Println("hello everybody , I'm asong")
 ch <- 1
}

func main()  {
 ch := make(chan int)
 go Hello(ch)
 <-ch
 fmt.Println("Golang梦工厂")
}
```

### 2. goroutine 的调度模型

- M代表线程
- P代表处理器，每一个运行的M（线程）都必须绑定一个P（处理器）
- G代表 goroutine，每次使用 go 关键字的时候，都会创建一个G对象

![GMP调度模型](/gmp.png "图3：GMP调度模型")

当前有两个P，各自绑定了一个M，每个P上挂了一个本地 goroutine 队列，也有一个全局 goroutine 队列。流程：

- 每次使用 go 关键字声明时，一个G对象被创建并加入到本地G队列或者全局G队列。
- 检查是否有空闲的P（处理器），若有那么创建一个M（若有正在 sleep 的M那么直接唤醒它）与其绑定，然后这个M循环执行 goroutine 任务。
- G任务执行的顺序是，先从本地队列中找。但若某个M（线程）发现本地队列为空，那么会从全局队列中截取 goroutine 来执行（一次性转移（全局队列的G个数/P个数））。如果全局队列也空，那么会随机从别的P那里截取 “一半” 的 goroutine 过来（偷窃任务），若所有的P的队列都为空，那么该M（线程）就会陷入 sleep。

### 3. goroutine 协程池

超大规模并发的场景下，不加限制的大规模的 goroutine 可能造成内存暴涨，给机器带来极大的压力，吞吐量下降和处理速度变慢还是其次，更危险的是可能使得程序 crash。所以，goroutine 池化是有其现实意义的。

{{< admonition type=question open=true >}}
首先，100w个任务，是不是真的需要100w个 goroutine 来处理？
{{< /admonition >}}

未必！用1w个 goroutine 也一样可以处理，让一个 goroutine 多处理几个任务就是了嘛，池化的核心优势就在于对 goroutine 的复用。此举首先极大减轻了 runtime 调度 goroutine 的压力，其次，便是降低了对内存的消耗。


{{< admonition type=tip open=true >}}
Goroutine Pool 的实现思路大致如下：

启动服务之时先初始化一个 Goroutine Pool 池，这个 Pool 维护了一个类似栈的 LIFO 队列 ，里面存放负责处理任务的 Worker，然后在 client 端提交 task 到 Pool 中之后，在 Pool 内部，接收 task 之后的核心操作是：

- 检查当前 Worker 队列中是否有空闲的 Worker，如果有，取出执行当前的 task；
- 没有空闲 Worker，判断当前在运行的 Worker 是否已超过该 Pool 的容量，是 — 阻塞等待直至有 Worker 被放回 Pool；否 — 新开一个 Worker（goroutine）处理；
- 每个 Worker 执行完任务之后，放回 Pool 的队列中等待。
{{< /admonition >}}

按照这个设计思路，一个高性能的 goroutine Pool，较好地解决了上述的大规模调度和资源占用的问题，在执行速度和内存占用方面相较于原生 goroutine 并发占有明显的优势，尤其是内存占用，因为复用，所以规避了无脑启动大规模 goroutine 的弊端，可以节省大量的内存。


## 3 等待组 WaitGroup

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

## 4 channel

`channel` 主要采用 CSP 并发模型实现的原理：不要通过共享内存来通信，而要通过通信来实现内存共享。它分为两种：带缓冲、不带缓冲。对不带缓冲的 `channel` 进行的操作实际上可以看作 “同步模式”，带缓冲的则称为 “异步模式”。

### 1. 非缓冲的 channel

无缓冲的通道只有当发送方和接收方都准备好时才会传送数据, 否则准备好的一方将会被阻塞。

### 2. 带缓冲的 channel

有缓冲的 `channel` 区别在于只有当缓冲区被填满时, 才会阻塞发送者, 只有当缓冲区为空时才会阻塞接受者。值得注意的是，

- 关闭 `channel` 以后仍然可以读取数据
- `for range` 循环可以持续从一个 `channel` 中接收数据

### 3. channel 的底层实现

#### 1 channel 底层结构体

![channel1](/channel1.png "图4：channel 底层结构体")

- `buf` 是有缓冲的 `channel` 所特有的结构，用来存储缓存数据。是个循环链表
- `sendx` 和 `recvx` 用于记录 `buf` 这个循环链表中的~发送或者接收的 `~index`
- `lock` 是个互斥锁。
- `recvq` 和 `sendq` 分别是接收 (`<-channel`) 或者发送 (`channel <- xxx`) 的 goroutine 抽象出来的结构体 (`sudog`) 的队列。是个双向链表

`channel` 的实现借助于结构体 `hchan`, 如下：
```go
type hchan struct {
    qcount   uint           // total data in the queue
    dataqsiz uint           // size of the circular queue
    buf      unsafe.Pointer // points to an array of dataqsiz elements
    elemsize uint16
    closed   uint32
    elemtype *_type // element type
    sendx    uint   // send index
    recvx    uint   // receive index
    recvq    waitq  // list of recv waiters
    sendq    waitq  // list of send waiters

    // lock protects all fields in hchan, as well as several
    // fields in sudogs blocked on this channel.
    //
    // Do not change another G's status while holding this lock
    // (in particular, do not ready a G), as this can deadlock
    // with stack shrinking.
    lock mutex
}
```

#### 2 send/recv 的细化操作

缓存链表中以上每一步的操作，都是需要加锁操作的！

每一步的操作的细节可以细化为：

- 第一，加锁
- 第二，把数据从 goroutine 中 copy 到“队列”中(或者从队列中 copy 到 goroutine 中）。
- 第三，释放锁

goroutine 内存 copy 到 `channel`:

![channel2](/channel2.png "图5：内存条 copy 进 channel")

`channel` 中的内存 copy 到 goroutine:

![channel3](/channel3.png "图6：channel 内存 copy 到内存条")

#### 3. goroutine 的阻塞操作

goroutine 的阻塞操作，实际上是调用 `send` (`ch <- xx`) 或者 `recv` ( `<-ch`) 的时候主动触发的，

```go
//goroutine1 中，记做G1

ch := make(chan int, 3)

ch <- 1
ch <- 1
ch <- 1
```

当 `channel` 缓存满了以后，再次进行 `send` 操作 (`ch<-1`) 的时候，会主动调用Go的调度器,让G1等待，并从让出M，让其他G去使用，

![channel4](/channel4.png "图7：Goroutine 调度")

同时G1也会被抽象成含有G1指针和 `send` 元素的 `sudog` 结构体保存到 `hchan` 的 `sendq` 中等待被唤醒。直到另一个 goroutine G2从缓存队列中取出数据，`channel` 会将等待队列中的G1推出，将G1当时 `send` 的数据推到缓存中，然后调用 Go 的 scheduler，唤醒G1，并把G1放到可运行的 goroutine 队列中。

### 4. channel 可能出现的状态

| 操作      | nil的channel | 正常的channel | 已关闭的channel |
| --------- | ------------ | ------------- | --------------- |
| <- ch     | 阻塞         | 成功或阻塞    | 读到零值        |
| ch <-     | 阻塞         | 成功或阻塞    | panic           |
| close(ch) | panic        | 成功          | panic           |

## 5 context 包

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

## 6 defer，panic 和 recover

### 1. defer
- 规则一：延迟函数的参数在 `defer` 语句出现时就已经确定下来了。例如：

  `defer` 语句中的 `fmt.Println()` 参数i值在 defer 出现时就已经确定下来，实际上是拷贝了一份。后面对变量i的修改不会影响 `fmt.Println()` 函数的执行，仍然打印"0"。

{{< admonition type=example >}}
```go
func a() {    
    i := 0
    defer fmt.Println(i)
    i++
    return
}
```
{{< /admonition >}}

- 规则二：延迟函数执行按后进先出顺序执行，即先出现的 `defer` 最后执行。

  定义 `defer` 类似于入栈操作，执行 `defer` 类似于出栈操作。设计 `defer` 的初衷是简化函数返回时资源清理的动作，资源往往有依赖顺序，比如先申请A资源，再跟据A资源申请B资源，跟据B资源申请C资源，即申请顺序是:A-->B-->C，释放时往往又要反向进行。这就是把 `deffer` 设计成FIFO的原因。每申请到一个用完需要释放的资源时，立即定义一个 `defer` 来释放资源是个很好的习惯。例如：

  函数拥有一个具名返回值 `result`，函数内部声明一个变量`i`，`defer` 指定一个延迟函数，最后返回变量i。延迟函数中递增 `result`。

  函数输出2。函数的 `return` 语句并不是原子的，实际执行分为设置返回值-->ret，`defer` 语句实际执行在返回前，即拥有 `defer` 的函数返回过程是：设置返回值-->执行 defer-->ret。所以 `return` 语句先把 `result` 设置为i的值，即1，`defer` 语句中又把 `result` 递增1，所以最终返回2。

{{< admonition type=example >}}
```go
func deferFuncReturn() (result int) {
  i := 1

  defer func() {
    result++
  }()

  return i
}

```
{{< /admonition >}}

- 规则三：延迟函数可能操作主函数的具名返回值。
   
  定义 `defer` 的函数，即主函数可能有返回值，返回值有没有名字没有关系，`defer` 所作用的函数，即延迟函数可能会影响到返回值。例如我们再看一下上面 `deferFuncReturn()` 的例子:

{{< admonition type=example >}}
```go
func deferFuncReturn() (result int) {
  i := 1

  defer func() {
    result++
  }()

  return i
}
```
{{< /admonition >}}

该函数的 `return` 语句可以拆分成下面两行：
```go
result = i
return
```
而延迟函数的执行正是在 `return` 之前，即加入 `defer` 后的执行过程如下：
```go
result = i
result++
return
```
所以返回值为 `result=1`。


### 2. panic

- 主动：程序猿主动调用 `panic()` 函数；
- 被动：编译器的隐藏代码触发，或者内核发送给进程信号触发；

`panic` 的具体实现，是依靠 `defer` 指针处理的，我们先来看一看 `panic` 的结构体：

```go
//runtime/runtime2.go
type _panic struct {
    argp unsafe.Pointer // pointer to arguments of deferred call run during panic; cannot move - known to liblink
    arg interface{} // argument to panic
    link *_panic // link to earlier panic
    recovered bool // whether this panic is over
    aborted bool // the panic was aborted
}
```
`_panic` 是个结构体，存储了 `defer` 指针、参数，`panic` 列表的表头指针，和已恢复或已终止的信息。以下是 `panic` 的处理流程：

- 1. 每个 goroutine 都有一个 `panic` 链表，运行时，遇到 `panic` 代码，会生成对应的 `_panic` 数据，存到这个链表的表头。
- 2. 每执行完毕一个函数，如果没有 `panic` 发生，就跳过对应的 `_panic` 数据，回到正常流程，否则进入3。
- 3. 如果有 `panic` 发生，处理链表中对应的 `_panic`，进入4。
- 4. 如果 `defer` 链表（跟 `panic` 链表一样，也是每个 goroutine 一个）里存在 `defer`，按约定顺序执行延迟代码，进入5，否则进入8。
- 5. 当 `defer` 链表执行到需要 `recover` 的时候，就交给 `reflectcall` 去调用 gorecover，进入6，否则进入7。
- 6. 执行 `recover`，这时对应的 `_panic` 结构里的 recovered 字段标记为真，由 recovery 方法，负责安抚当前的 `_panic`，回到正常流程。
- 7. 如果没 `recover`，那就进入死给你看流程，进入8。
- 8. 最后，执行 `fatalpanic` 方法。

注意：因为 golang 的 gorotuine 机制，`panic` 在不同的 gorotuine 里面，是单独的，并不是整体处理。可能一个地方凉了，就会整体完蛋，这个要非常小心。

### 3. recover

Golang 虽然没有 `try catch` 机制，但它有类似 `recover` 的机制，

{{< admonition type=example >}}
```go
package main

import "fmt"

func main() {
	fmt.Printf("%d\n", cal(1, 2))
	fmt.Printf("%d\n", cal(5, 2))
	fmt.Printf("%d\n", cal(5, 0))
	fmt.Printf("%d\n", cal(9, 2))
}

func cal(a, b int) int {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("%s\n", err)
		}
	}()
	return a / b
}
```
{{< /admonition >}}

在 `cal` 函数里面每次终止的时候都会检查有没有异常产生，如果产生了我们可以处理，比如说记录日志，这样程序还可以继续执行下去。

## 7 interface 

### 1. 多态

`interface` 定义了一个或一组 method(s)，若某个数据类型实现了 `interface` 中定义的那些被称为 "methods" 的函数，则称这些数据类型实现（implement）了 `interface`。举个例子来说明。

{{< admonition type=example >}}
```go
//定义了一个Mammal的接口，当中声明了一个Say函数。
type Mammal interface {
 Say()
}
```
{{< /admonition >}}

定义 `Cat`、`Dog` 和 `Human` 三个结构体，分别实现各自的 `Say` 方法：

```go
type Dog struct{}

type Cat struct{}

type Human struct{}

func (d Dog) Say() {
 fmt.Println("woof")
}

func (c Cat) Say() {
 fmt.Println("meow")
}

func (h Human) Say() {
 fmt.Println("speak")
}
```

之后，我们尝试使用这个接口来接收各种结构体的对象，然后调用它们的 `Say` 方法：

```go
func main() {
 var m Mammal
 m = Dog{}
 m.Say()
 m = Cat{}
 m.Say()
 m = Human{}
 m.Say()
}
// print result:
// woof
// meow
// speak
```
### 2. 类型断言 (type assertion)

```go
//类型断言
//一个判断传入参数类型的函数
func just(items ...interface{}) {
    for index, v := range items {
        switch v.(type) {
        case bool:
            fmt.Printf("%d params is bool,value is %v\n", index, v)
        case int, int64, int32:
            fmt.Printf("%d params is int,value is %v\n", index, v)
        case float32, float64:
            fmt.Printf("%d params is float,value is %v\n", index, v)
        case string:
            fmt.Printf("%d params is string,value is %v\n", index, v)
        case Student:
            fmt.Printf("%d params student,value is %v\n", index, v)
        case *Student:
            fmt.Printf("%d params *student,value is %v\n", index, v)

        }
    }
}
```
我们聊了面向对象中多态以及接口、类型断言的概念和写法，借此进一步了解了为什么 golang 中的接口设计非常出色，因为它**解耦了接口和实现类之间的联系**，使得进一步增加了我们编码的灵活度，解决了供需关系颠倒的问题。

## 8 hashmap

{{< admonition type=question open=true >}}
Hashmap 的内部结构是如何实现的呢？
{{< /admonition >}}

### 1. 内部结构

Go的 map 是 unordered map，即无法对 key 值排序遍历。跟传统的 hashmap 的实现方法一样，它通过一个 buckets 数组实现，所有元素被 hash 到数组的 bucket 中，buckets 就是指向了这个内存连续分配的数组。`B`字段说明 hash 表大小是2的指数，即$2^B$。每次扩容会增加到上次大小的两倍，即 $2^{B+1}$。当 bucket 填满后，将通过 `overflow` 指针来 `mallocgc` 一个bucket出来形成链表，也就是为哈希表解决冲突问题。

```go
// A header for a Go map.
type hmap struct {
	count int // len()返回的map的大小 即有多少kv对
	flags uint8
	B     uint8  // 表示hash table总共有2^B个buckets 
	hash0 uint32 // hash seed
	buckets    unsafe.Pointer // 按照low hash值可查找的连续分配的数组，初始时为16个Buckets.
	oldbuckets unsafe.Pointer 
	nevacuate  uintptr      
	overflow *[2]*[]*bmap //溢出链 当初始buckets都满了之后会使用overflow
}

// A bucket for a Go map.
type bmap struct {
    tophash [bucketCnt]uint8
// Followed by bucketCnt keys and then bucketCnt values.
// NOTE: packing all the keys together and then all the values together makes the
// code a bit more complicated than alternating key/value/key/value/... but it allows
// us to eliminate padding which would be needed for, e.g., map[int64]int8.
// Followed by an overflow pointer.
}
```
bucket map 的数据结构中，`tophash` 是个大小为 8(bucketCnt) 的数组，存储了8个 key 的 hash 值的高八位值。

{{< admonition type=tip open=true >}}
在对 key/value 对增删查的时候，先比较 key 的 hash 值高八位是否相等，然后再比较具体的key值。
{{< /admonition  >}}

### 2. map 初始化
### 3. hash grow 扩容和迁移


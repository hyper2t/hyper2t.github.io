# Go 基础知识与框架体系 part 3: channel


这篇文章总结了 Go 的知识体系: channel，包括其中的底层实现等等。

<!--more-->


## channel

`channel` 主要采用 CSP 并发模型实现的原理：不要通过共享内存来通信，而要通过通信来实现内存共享。它分为两种：带缓冲、不带缓冲。对不带缓冲的 `channel` 进行的操作实际上可以看作 “同步模式”，带缓冲的则称为 “异步模式”。

### 1. 非缓冲的 channel

无缓冲的通道只有当发送方和接收方都准备好时才会传送数据, 否则准备好的一方将会被阻塞。

### 2. 带缓冲的 channel

有缓冲的 `channel` 区别在于只有当缓冲区被填满时, 才会阻塞发送者, 只有当缓冲区为空时才会阻塞接受者。值得注意的是，

- 关闭 `channel` 以后仍然可以读取数据
- `for range` 循环可以持续从一个 `channel` 中接收数据

### 3. channel 的底层实现

#### 1 channel 底层结构体

![channel1](/channel1.png "图1：channel 底层结构体")

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

![channel2](/channel2.png "图2：内存条 copy 进 channel")

`channel` 中的内存 copy 到 goroutine:

![channel3](/channel3.png "图3：channel 内存 copy 到内存条")

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

![channel4](/channel4.png "图4：Goroutine 调度")

同时G1也会被抽象成含有G1指针和 `send` 元素的 `sudog` 结构体保存到 `hchan` 的 `sendq` 中等待被唤醒。直到另一个 goroutine G2从缓存队列中取出数据，`channel` 会将等待队列中的G1推出，将G1当时 `send` 的数据推到缓存中，然后调用 Go 的 scheduler，唤醒G1，并把G1放到可运行的 goroutine 队列中。

### 4. channel 可能出现的状态

| 操作      | nil的channel | 正常的channel | 已关闭的channel |
| --------- | ------------ | ------------- | --------------- |
| <- ch     | 阻塞         | 成功或阻塞    | 读到零值        |
| ch <-     | 阻塞         | 成功或阻塞    | panic           |
| close(ch) | panic        | 成功          | panic           |


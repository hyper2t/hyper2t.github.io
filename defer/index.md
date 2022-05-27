# Golang 基础知识与框架体系 part 6: defer, panic 和 recover


这篇文章总结了 Golang 的知识体系 `defer`, `panic` 和 `recover`，包括其中的底层实现等等。

<!--more-->

## defer，panic 和 recover

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


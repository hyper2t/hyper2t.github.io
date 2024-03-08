# Go 基础知识与框架体系 系列七: interface


这篇文章总结了 Golang 的知识体系: interface，包括其中的底层实现等等。

<!--more-->

{{< admonition >}}
说起 Golang， 大家都会第一时间想到高并发和 Golang 作为主流的后端开发语言的优势，本文主要讲 Golang 主要知识体系，包括**数组和切片**、**协程的调度**原理、等待组 **waitGroup**、**channel** 的底层实现、互斥锁 **mutex** 的实现、**interface** 中的多态等等。
{{< /admonition >}}

## interface

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


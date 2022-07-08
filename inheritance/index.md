# C++ 面向对象编程知识体系 part 2: 类的继承


这篇文章总结了 C++ 中面向对象编程的知识体系：类的继承，包括其中的底层实现等等。

<!--more-->

## 1 类的继承

### 1.1 继承的定义


### 1.2 虚继承

{{< admonition definition>}}
基类：
```cpp
class Parent {
// ...
};
```
子类：
```cpp
class Child : virtual public Parent {
// ...
};
```
{{< /admonition>}}

证明这个 `Parent` 是 `Child` 的虚基类。

虚基类是指：`class SubClass : virtual public BaseClass` 中以 `virtual` 声明的基类，由于 C++ 支持多重继承，所以对于一个派生类中有几个直接父类，而几个直接父类中有几个可能分别继承自某一个基类（就是父类的父类），这样在构造最终派生类时，会出现最终派生类中含有多个同一个基类的情况，就会产生二义性的问题（不知道该调用哪个基类的成员变量和函数），为解决此问题，需要使用虚基类，即只对此基类生成一块内存区域，这样最终派生类中就只会含有一个基类了。

典型的需要用虚基类的情况如下(菱形继承)：
```cpp
    A
   / \
  B   C
   \ /
    D
```
其中 `D` 继承自 `B C`，`B C` 分别继承自 `A`，所以 `A` 要分别被 `B C` 虚拟继承：

```cpp
#include <iostream>

using namespace std;

class A {
  public:
   void printA() { cout << "this is A\n"; }
};

class B : virtual public A {};
class C : virtual public A {};
class D : public B, public C {};

int main(int argc, char const* argv[]) {
  D* d = new D;
  d->printA();
  return 0;
}
```

如果把 virtual 虚继承去掉，会报错为：

```cpp
error: request for member ‘printA’ is ambiguous
```

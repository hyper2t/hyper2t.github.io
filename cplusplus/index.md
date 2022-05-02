# C++ 面向对象编程知识体系


这篇文章总结了 C++ 中面向对象编程的知识体系和实际应用中用到的框架，包括其中的底层实现等等。

<!--more-->

## 1 面向对象特性

### 1. 类与对象

{{< admonition question>}}
什么是**类与对象**？
{{< /admonition >}}

- **对象**是现实中的对象在程序中的模拟
- **类**是同一类对象的抽象，**对象**是类的某一特定实体
- **类**是一种用户自定义的类型，包含函数与数据的特殊结构体

{{< admonition example>}}
举一个简单的例子，
```c++
class Orange{
public:
  int GetOrangeDegree(); // orange 的成员方法
  int GetSourDegree();
  int GetSweetDegree();
  int GetWaterRatio();
private:
  int m_waterRatio;
  int m_orangeDegree;
  int m_sourDegree;
  int m_sweetDegree;
}
```
{{< /admonition >}}

这里所举的例子，`Orange`是类，但橙子有分不同个体，它们都有对应的酸度、甜度、水分和生长高度，所以不同的橙子个体代表着这一类橙子的不同对象。

#### **类成员的访问控制**
- 访问权限分为公有类型`public`，保护类型`protected`，私有类型`private`
- 成员默认访问权限为`private`
- 友元函数或者友元类可访问类的保护成员或私有成员

来看一个例子：
{{< admonition example>}}
```c++
class Teacher;
class Student{
public:
  void SetId(uint32_t id){m_id = id;}
  uint32_t GetId();
  friend Teacher;
  friend void PrintStudentId(Student);
private:
  uint32_t m_id = 0;
  uint32_t m_age = 18;
}

class Teacher{
public:
  void PrintStudentId(Student stud)
  {
    std::cout << stud.m_id << std:endl;
  }
}

void PrintStudentId(Student stud)
{
  std::cout << stud.m_id << std:endl;
}

```
{{< /admonition >}}

我们定义了 `class Teacher`, `Teacher`是`Student`的友元类，`PrintStudentId`是`Student`的友元函数，该函数可用于访问`Student`类的保护成员或者私有成员。

例如：就像上面所写的，`Teacher`类里面定义的`PrintStudentId`函数可以访问`Student`类的私有成员`stud.m_id`。

{{< admonition warning>}}
1. 友元函数是**单向**的： 打个比方，`Student`把`Teacher`当成朋友，而`Teacher`不把`Student`当朋友。就如上面的例子，`Teacher`是`Student`的友元，但是`Student`不是`Teacher`的友元，所以`Teacher`类的成员函数在没有`friend`的关键字作用下，不能访问`Teacher`类的保护成员或者私有成员。
2. 友元函数不是类的成员。
{{< /admonition>}}

### 2. 类的构造与析构

#### **构造函数**
- **构造函数**是用于构造函数的特殊函数，在对象被创建时被调用以初始化对象
- 未定义**构造函数**时，编译器自动生成不带参数的默认版本
- 执行**构造函数**时先执行其初始化列表，再执行函数体
- **构造函数**和其他函数一样，允许被重载和被委托
- **构造函数**的名称与类的名称是完全相同的，并且不会返回任何类型，也不会返回`void`。构造函数可用于为某些成员变量设置初始值。

**构造函数**主要有以下三个方面的作用：
- 给创建的对象建立一个标识符；
- 为对象数据成员开辟内存空间；
- 完成对象数据成员的初始化。
```c++
class Student{
public:
  Student();
  Student(uint32_t id);
  Student(uint32_t id, uint32_t age);
  void SetId(uint32_t id){m_id = id;}
  uint32_t GetId(){return m_id;}
private:
  uint32_t m_id = 0;
  uint32_t m_age;
}
```
关于以上的`Student`类，我们声明了三个`Student`的构造函数。第一个构造函数`Student()`不带参数。三个构造函数的具体定义如下：

```c++
Student::Student(uint32_t id, uint32_t age):m_id(id), m_age(age)
{
  // do initialization
}
Student::Student():Student(0,18)
{
  // do initialization
}
Student::Student(uint32_t id):Student(id,18)
{
  // do initialization
}
```

按照上面的写法，我们可以使用初始化列表进行初始化字段。假设有一个类`C`，具有多个字段`X、Y、Z`等需要进行初始化，同理地，您可以使用上面的语法，只需要在不同的字段使用逗号进行分隔，如下所示：
```c++
C::C( double a, double b, double c): X(a), Y(b), Z(c)
{
  ....
}
```

{{< admonition question>}}
既然构造函数允许被重载，构造函数是如何做到被重载呢？
{{< /admonition>}}

需要注意的是, 在进行构造函数的重载时要注意重载和参数默认的关系要处理好, 避免产生代码的二义性导致编译出错, 例如以下具有二义性的重载:

{{< admonition example>}}
```c++
Point(int x = 0, int y = 0)     //默认参数的构造函数
{
  xPos = x;
  yPos = y;
}

Point()     //重载一个无参构造函数
{
  xPos = 0;
  yPos = 0;
}
```
{{< /admonition>}}

在上面的重载中, 当尝试用`Point`类重载一个无参数传入的对象`M`时, `Point M`; 这时编译器就报一条`error: call of overloaded 'Point()' is ambiguous`的错误信息来告诉我们说 `Point` 函数具有二义性。

这是因为`Point(int x = 0, int y = 0)`全部使用了默认参数, 即使我们不传入参数也不会出现错误, 但是在重载时又重载了一个不需要传入参数了构造函数`Point()`, 这样就造成了当创建对象都不传入参数时编译器就不知道到底该使用哪个构造函数了, 就造成了二义性。

#### **析构函数**

与构造函数相反, **析构函数**是在对象被撤销时被自动调用, 用于对成员撤销时的一些清理工作, 例如在前面提到的手动释放使用`new`或`malloc`进行申请的内存空间。**析构函数**具有以下特点:
- **析构函数**函数名与类名相同, 紧贴在名称前面用波浪号 ~ 与构造函数进行区分, 例如: `~Point()`;
- 构造函数没有返回类型, 也不能指定参数, 因此**析构函数**只能有一个, 不能被重载;
- 当对象被撤销时析构函数被自动调用, 与构造函数不同的是, **析构函数**可以被显式的调用, 以释放对象中动态申请的内存。

某种意义上理解析构函数是回收间接资源的函数。

![析构函数](/析构函数.png "图1：析构函数回收间接资源")

```c++
#include <iostream>
#include <cstring>
using namespace std;

class Book
{
 public:
     Book( const char *name )      //构造函数
    {
        bookName = new char[strlen(name)+1];
        strcpy(bookName, name);
    }
    ~Book()                 //析构函数
    {
        cout<<"析构函数被调用...\n";
        delete []bookName;  //释放通过new申请的空间
    }
    void showName() { cout<<"Book name: "<< bookName <<endl; }

 private:
    char *bookName;
};

int main()
{
    Book CPP("C++ Primer");
    CPP.showName();

    return 0;

}
```

代码中创建了一个`Book`类, 类的数据成员只有一个字符指针型的`bookName`, 在创建对象时系统会为该指针变量分配它所需内存, 但是此时该指针并没有被初始化所以不会再为其分配其他多余的内存单元。在构造函数中, 我们使用`new`申请了一块`strlen(name)+1`大小的空间, 也就是比传入进来的字符串长度多1的空间, 目的是让字符指针`bookName`指向它, 这样才能正常保存传入的字符串。

在`main`函数中使用`Book`类创建了一个对象`CPP`, 初始化`bookName`属性为`"C++ Primer"`。从运行结果可以看到, 析构函数被调用了, 这时使用`new`所申请的空间就会被正常释放。

自然状态下对象何时将被销毁取决于对象的生存周期, 例如全局对象是在程序运行结束时被销毁, 自动对象是在离开其作用域时被销毁。

如果需要显式调用析构函数来释放对象中动态申请的空间只需要使用 对象名.析构函数名(); 即可, 例如上例中要显式调用析构函数来释放`bookName`所指向的空间只要:

```c++
CPP.~Book();
```
### 3. 类的其他特性

### 4. 类的继承

### 5. 类的多态

### 6. RTTI 与抽象类

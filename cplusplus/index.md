# C++ 基础知识与框架体系


这篇文章总结了 C++ 的知识体系和实际应用中用到的框架，包括其中的底层实现等等。

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

**类成员的访问控制**
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

**关于构造函数**：
- **构造函数**是用于构造函数的特殊函数，在对象被创建时被调用以初始化对象
- 未定义**构造函数**时，编译器自动生成不带参数的默认版本
- 执行**构造函数**时先执行其初始化列表，再执行函数体
- **构造函数**和其他函数一样，允许被重载和被委托

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

### 3. 类的其他特性

### 4. 类的继承

### 5. 类的多态

### 6. RTTI 与抽象类

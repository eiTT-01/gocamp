# 1 Go语言的变量命名规范是什么？

首字母大写为公开函数或变量，可以被其他包调用。
否则为私有变量，只能在包内调用。

# 如何获取变量的地址？如何取到地址中的值？

获取变量地址: `&p`
读取指针指向的变量的值: `*p`

# 变量的生命周期是什么？作用域是什么？

对于包一级声明的变量，它的生命周期和整个程序的运行周期是一致的。而局部变量的生命周期相对动态：从它被创建开始，直到该变量不在被引用为止。

# 创建变量有哪几种方式

```
var a int
b := 1
var c = 1
```

# Go语言简单数据类型都有哪些？

```
// 整型
var i int = 10
// 浮点数
var f float32 = 1.234
// 复数
var x complex128 = complex(1, 2)
x := 1 +2i
// 布尔型
var b bool = true
// 字符串
s := "Hello Go"
// 常量
const pi = 3.1415
```

# 初始化数组的几种方式？

```
var a [3]int
b := [...]int{1, 2, 3}
c := [3]int{1, 2, 3}
```

# 遍历数组

```
arr := [3]int{1, 2, 3}

for _, i := range arr {
    fmt.Println(i)
}
```


# 初始化切片的几种方式？

```
// slice
var s []int
s := make([]int, 5, 10)
s := make([]int, 5)
s := []int{1, 2, 3}
```

# 如何复制切片

```
s1 := []int{1, 2, 3}

// 浅复制
s2 := s1
// 深复制
s3 := make([]int, 3)
copy(s3, s2)
```

# 实现切片的增删改查

```
s := []int{1,2,3,4}

// 增
s = append(s, 5)
// 删
s = s[:len(s)-1]
// 改
s[1] = 5
// 查
for _, item := ragne {
    fmt.Println(item)
} 

```

# 下面代码是否有问题？并说出为什么？ 如果有问题，如何修正？

```
s := []string{"炭烤生蚝", "麻辣小龙虾", "干锅鸭"}
s2 := make([]*string, len(s))

for i, v := range s {
    s2[i] = &s[i]
}
```

上述例子中，循环体的就变量 `v` 是值的副本，每次遍历 `s` 时都是一次值拷贝，它指向的地址是不变的。也就是说，`s2` 中的每个元素都保存了 `v` 的地址。遍历结束后，`v` 的地址将指向的变量值为 `s` 的最后一个元素，即"干锅鸭"。

修改之后的程序，

```
s := []string{"炭烤生蚝", "麻辣小龙虾", "干锅鸭"}
s2 := make([]*string, len(s))

for i, v := range s {
	val := v // 在局部作用域内拷贝一份
    s2[i] = &val
}
```


# 分别写一个 if 和 switch、枚举 的例子

```
var a int = 10

if a < 100 {
	fmt.Println("a < 100")
} else {
	fmt.Println("a >= 100")
}

switch a {
    case 1:
		fmt.Println("a == 1")
    case 10:
		fmt.Println("a == 10")
    default:
		fmt.Println("a != 1 and 10")
}

```

# map有什么特点？

map 就是一个哈希表的引用。map 中的所有的 key 都是相同的类型，所有的 value 也都是相同的类型，但是 key 和 value 可以是不同类型。另外，在哈希表中，key 是不允许重复的，而且通过 key 去查找 value 的时间复杂度为 1，效率非常高。


# 什么样的类型可以做map的key

除了 slice、map 和 function 以外，基本都可以作为 map 的 key。


# 写一个map的增删改查

```
m := map[string]int{
	"alice": 10,
	"bob": 11,
}
// 增
m["charlie"] = 12
// 删
delete(m, "charlie")
// 改
m["alice"] = 13
// 查
if bob, ok := m["bob"]; ok {
	fmt.Println(bob)
}
```

# 函数的定义

```
func greet(name string) {
	fmt.Printf("Hello %s!", name)
}
```


# 函数传参，传值还是传引用？

在 Go 语言中，所有的传参都是值传递。

# 定义函数的多返回值？

```
func multiret(a, b int) (int, error) {
	return a + b, nil
}
```


# 举例说明 函数变量、匿名函数、闭包、变长函数？

```
// inputs 为可变参
func counter(inputs ...int) func int {
    x := 0 // 函数变量
    for _, v := range inputs {
        x += v
    }
    return func() int { // 匿名函数
        x++  // 闭包
        return x
    }
}
```

# 说一下面向对象设计的好处？

面向对象的设计有诸多好处，其中包括但不限于易维护和易扩展。

易维护：可读性高，在继承的帮助下，即使改变需求，也可以将改动限制在局部模块。
易扩展：由于继承、封装、多态等特性，很容易设计出高内聚低耦合的的系统结构，使得系统更灵活、更容易扩展。

# 方法的定义

```
type People struct{}

func (p People) SayHi() {
    fmt.Println("Hi")
}

```

# 指针接收者和值接收者有何不同

实现了接收者是值类型的方法，相当于自动实现了接收者是指针类型的方法；而实现了接收者是指针类型的方法，不会自动生成对应接收者是值类型的方法。

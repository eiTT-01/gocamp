# 说明一下接口是什么，和面向对象有什么关系？
（选做部分:如果你知道java，那么，Go语言的接口和java接口有什么不同？）

接口类型是对其他类型行为的抽象和概括。在 golang 中，接口就是一组方法签名。

# 举例说明鸭子类型

duck typing 描述事物的外部行为而非内部结构。简单来说，谁实现了接口A，谁就是A。

```
type Duck interface {
    Gaga()
}

type RealDuck struct {}
func (rd RealDuck) Gaga() {
    fmt.Println("Gaga")
}

type AlienDuck struct {}
func (ad AlienDuck) Gaga() {
    fmt.Println("Gaga")
}

func SayGaga(d Duck) {
    d.Gaga()
}

```


# go语言中的标准的接口，有哪些？ 并举例说明1-2个接口的实现，通过接口如何实现多态？

常见的有 `io.Writer`、`io.Reader`。

```
type ByteCounter int

func (ni *ByteCounter) Write(p []byte) (int ,error) {
    *c += ByteCounter(len(p))
    return len(p), nil
}
```

# 函数传值和传引用有何不同？ 各举一个例子


# 延长函数的调用顺序是什么？ 举例说明

defer 的执行顺序与声明顺序相反。这是因为 defer 声明的函数会被保存到一个类似栈的数据结构中，具有后进先出的特性。

```
function DeferTest() {
    defer fmt.Println("the last one to be executed")
    defer fmt.Println("the first one to be executed")
    
    fmt.Println("to be executed before defer func")
}

```

# go语言是如何做测试的？ 举例说明

通过 `go test` 命令执行测试文件，它是一个按照一定约定和组织来测试代码的命令。`go test` 命令会遍历所有 `_test.go` 文件中符合命名规范的函数。

```
func FuncA(input int) bool {
    return input < 2
}

func TestFuncA(t *testing.T) {
    var cases = []struct {
        input int
        want bool
    }{
        {"input": 1, want: true},
        {"input": 2, want: false},
    }
    for _, case := range cases {
        if ret := FuncA(case.input); !ret {
            t.Errorf("when input=%d, expect %t but got %t", case.input, case.want, ret)
        }
    }
}
```

# 如何理解线程安全？

当多线程程序需要访问共享变量/数据的时候，所有的线程都可以正常且正确地执行，不会出现数据污染的情况。要想保证线程安全很简单，就是让所有的线程能够串型依次去操作共享数据。

# 如何理解Go语言的并发模型？

在 Go 语言中，每一个并发的执行单元叫做一个 `goroutine`。Go 实现了两种并发形式。第一种就是常规的多线程共享内存，另一种是 Go 语言特有的 CSP 并发模型。

Go 的 CSP 并发模型，是通过 `goroutine` 和 `channel` 来实现的。


# 缓冲通道与无缓冲通道有何不同？

- 无缓冲通道会阻塞程序，直到有 `goroutine` 从通道中消费掉该数据为止。 
- 有缓冲的通道，在缓冲通道满之前，是不会阻塞程序执行的，当缓冲满之后，跟无缓冲通道一样，需要等待 `goroutine` 消费数据。

# 单向通道优势是什么？

可以在编译阶段就发现问题，防止 `channel` 的滥用。

# 关闭通道，会造成哪些影响？

- 任何写操作都将引发 `panic`
- 读一个已经关闭的通道，不会引发阻塞
- 对于一个关闭的通道，消费者依然可以读到之前写的数据，但当通道中已经没有数据时，将返回一个零值的数据（根据通道定义的数据类型）。
- 关闭一个只读的单向通道将引发编译错误

# 什么场景使用select?

当有多个通道需要监听时，为了避免程序被阻塞，就需要用到 `select` 语句。

# 举例说明 mutex 和 rwmutex

- Mutex 是互斥锁。互斥锁保证共享区域同时只能有一个 `goroutine` 访问和操作。
- RWMutex 是读写锁。它允许多个只读操作并行，但写操作完全互斥。

```
var mu sync.RWMutex
var balance int

func Balance() int {
    mu.RLock() // readers lock
    defer mu.RUnlock()
    return balance
}

func Deposit(amount int) {
    mu.Lock()
    defer mu.Unlock()
    balance += amount
}  

```

# 举例说明 条件变量

```
cond := sync.NewCond(new(sync.Mutex))
condition := 0

// 消费者
go func() {
    for {
        // 开始消费前，锁住
        cond.L.Lock()
        // 若没有可消费数据，等待
        for condition == 0 {
            cond.Wait()
        }
        // 消费
        condition--
        fmt.Printf("Consumer: %d\n", condition)
        
        // 唤醒生产者
        cond.Signal()
        // 解锁
        cond.L.Unlock()
    }
}()

// 生产者
for {
    // 生产前，锁住
    cond.L.Lock()
    // 生产太多时，等待消费
    for condition == 100 {
        cond.Wait()
    }
    // 生产
    condition++
    fmt.Printf("Producer: %d\n", condition)
    
    // 通知消费者
    cond.Signal()
    // 解锁
    cond.L.Unlock()
}
```

# 举例说明 WaitGroup

```
var wg sync.WaitGroup

wg.Add(1)
go func() { wg.Done() }()
wg.Wait()

```

# 举例说明 context.Context

上下文 `Context` 包提供了一组在 `goroutine` 间进行值传递的方法。

```
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	go handle(ctx, 500*time.Millisecond)

	select {
	case <-ctx.Done():
		fmt.Println("handle done in main", ctx.Err())
	}
}

func handle(ctx context.Context, d time.Duration) {
	select {
	case <-ctx.Done():
		fmt.Println("handle done in handle", ctx.Err())
	case <-time.After(d):
		fmt.Println("timeout in", d)
	}
}

```

# 说说你对GO语言错误处理的理解？

错误指的是可能出现问题的地方出现了问题，比如打开一个文件时失败，这种情况在人们的意料之中 ；而异常指的是不应该出现问题的地方出现了问题，比如引用了空指针，这种情况在人们的意料之外。可见，错误是业务过程的一部分，而异常不是。

go 语言中 error 接口类型作为错误处理的标准模式，如果函数要返回错误，则返回值类型列表中肯定包含 error。error 处理过程类似于 C 语言中的错误码，可逐层返回，直到被处理。

# go语言如何做依赖管理？

通过 `Go Modules` 来进行依赖管理。

# go mod 常用命令有哪些？

```shell
# 初始化 go mod
$ go mod init
# 安装依赖包
$ go get go.uber.org/zap
# 升级依赖包
$ go get -u go.uber.org/zap
# 整理依赖关系
$ go mod tidy
```

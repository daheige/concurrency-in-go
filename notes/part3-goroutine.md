# goroutine
Goroutine是Go中最基本的组织单位之一，所以了解它是什么以及它如何工作是非常重要的。`事实上，每个Go程序至少拥有一个：main gotoutine，当程序开始时会自动创建并启动。在几乎所有Go程序中，你都可能会发现自己迟早加入到一个gotoutine中，以帮助自己解决问题。`那么它到底是什么？
简单来说，gotoutine是一个并发的函数（记住：不一定是并行）和其他代码一起运行。你可以简单的通过将go关键字放在函数前面来启动它：
```go
func main() {
    go sayHello()
    // continue doing other things
}

func sayHello() {
    fmt.Println("hello")
}
```

对于匿名函数，同样也能这么干，从下面这个例子你可以看得很明白。在下面的例子中，我们不是从一个函数建立一个goroutine，而是从一个匿名函数创建一个goroutine：
```go
go func() {
    fmt.Println("hello")
}() // 1
// continue doing other things
```

注意这里的()，我们必须立刻调用匿名函数来使go关键字有效。
或者，你可以将函数分配给一个变量，并像这样调用它：

sayHello := func() {
    fmt.Println("hello")
}
go sayHello()
// continue doing other things
看起来很简单，对吧。我们可以用一个函数和一个关键字创建一个并发逻辑块，这就是启动goroutine所需要知道的全部。当然，关于如何正确使用它，对它进行同步以及如何组织它还有很多需要说明的内容。
本章接下来的部分会深入介绍goroutine及它是如何工作的。如果你只想编写一些可以在goroutine中正确运行的代码，那么可以考虑直接跳到下一章。
那么让我们来看看发生在幕后的事情：goroutine实际上是如何工作的？ 是OS线程吗？ 绿色线程？ 我们可以创建多少个？
`Goroutines对Go来说是独一无二的（尽管其他一些语言有类似的并发原语）。它们不是操作系统线程，它们不完全是绿色的线程(由语言运行时管理的线程)，它们是更高级别的抽象，被称为协程(coroutines)。协程是非抢占的并发子程序，也就是说，它们不能被中断。`

`Go的独特之处在于goutine与Go的运行时深度整合。Goroutine没有定义自己的暂停或再入点; Go的运行时观察着goroutine的行为，并在阻塞时自动挂起它们，然后在它们变畅通时恢复它们。在某种程度上，这使得它们可以抢占，但只是在goroutine被阻止的地方。它是运行时和goroutine逻辑之间的一种优雅合作关系。 因此，goroutine可以被认为是一种特殊的协程。`

`协程，因此可以被认为是goroutine的隐式并发构造，但并发并非协程自带的属性：某些东西必须能够同时托管几个协程，并给每个协程执行的机会，否则它们无法实现并发。`当然，有可能有几个协程按顺序执行，但看起来就像并行一样，在Go中这样的情况比较常见。

Go的宿主机制实现了所谓的M：`N调度器，这意味着它将M个绿色线程映射到N个系统线程。 Goroutines随后被安排在绿色线程上。 当我们拥有比绿色线程更多的goroutine时，调度程序处理可用线程间goroutines的分布，并确保当这些goroutine被阻塞时，可以运行其他goroutines。我们将在第六章讨论所有这些机制是如何工作的，但在这里我们将介绍Go如何对并发进行建模。`

`Go遵循称为fork-join模型的并发模型.fork这个词指的是在程序中的任何一点，它都可以将一个子执行的分支分离出来，以便与其父代同时运行。join这个词指的是这样一个事实，即在将来的某个时候，这些并发的执行分支将重新组合在一起。子分支重新加入的地方称为连接点。`
go关键字为Go程序实现了fork，fork的执行者是goroutine，让我们回到之前的例子：
```go
sayHello := func() {
    fmt.Println("hello")
}
go sayHello()
// continue doing other things
```

sayHello函数会在属于它的goroutine上运行，与此同时程序的其他部分继续执行。在这个例子中，没有连接点。执行sayHello的goroutine将在未来某个不确定的时间退出，并且该程序的其余部分将继续执行。

然而，这个例子存在一个问题：我们不确定sayHello函数是否可以运行。goroutine将被创建并交由Go的运行时安排执行，但在main goroutine退出前它实际上可能没有机会运行。

事实上，由于我们为了简单而省略了其他主要功能部分，所以当我们运行这个小例子时，几乎可以肯定的是，程序将在主办sayHello调用的goroutine开始之前完成执行。 因此，你不会看到打印到标准输出的单词“hello”。 你可以在创建goroutine之后为main goroutine添加一段休眠时间，但请记住，这实际上并不创建一个连接点，只是一个竞争条件。如果你记得第一章，你会增加退出前goroutine将运行的可能性，但你无法保证它。加入连接点是确保程序正确性并消除竞争条件的保证。

为了创建一个连接点，你必须同步main goroutine和sayHello goroutine。 这可以通过多种方式完成，但我将使用sync包中提供的一个解决方案：sync.WaitGroup。现在了解这个示例如何创建一个连接点并不重要，只是需要清楚它在两个goroutine之间创建了一个连接点。 这是我们的示例版本：
```
var wg sync.WaitGroup
sayHello := func() {
    defer wg.Done()
    fmt.Println("hello")
}
wg.Add(1)
go sayHello()
wg.Wait() //1
```

在这里加入连接点。
这会输出：
```
hello
```

这个例子明确的阻塞了main goroutine，直到承载sayHello函数的main goroutine终止。你将在随后的sync包章节了解到更详细的内容。

我们在示例中使用了匿名函数。让我们把注意力转移到闭包。闭包围绕它们创建的词法范围，从而捕捉变量。如果在goroutine中使用闭包，闭包是否在这些变量或原始引用的副本上运行？让我们试试看：
```go
var wg sync.WaitGroup
salutation := "hello"
wg.Add(1)
go func() {
    defer wg.Done()
    salutation = "welcome" // 1
}()
wg.Wait()
fmt.Println(salutation)
// 你认为salutation的值是"hello"还是"welcome"？运行后会看到：

wlecome
```

有趣！事实证明，goroutine在它创建的同一地址空间内执行，因此我们的程序打印出“welcome”。让我们再来尝试一个例子。 你认为这个程序会输出什么？
```go
var wg sync.WaitGroup
for _, salutation := range []string{"hello", "greetings", "good day"} {
    wg.Add(1)
    go func() {
        defer wg.Done()
        fmt.Println(salutation) // 1
    }()
}
wg.Wait()
```

这里我们测试打印字符串切片创建的循环变量salutation。
答案比大多数人所预期的不同，而且是Go中为数不多的令人惊讶的事情之一。 大多数人直觉上认为这会以某种不确定的顺序打印出“hello”，”greeting”和“good day”，但实际上：
```
good day 
good day 
good day
```

这有点令人惊讶。让我们来看看这里发生了什么。 在这个例子中，goroutine正在运行一个已经关闭迭代变量salutation的闭包，它有一个字符串类型。 当我们的循环迭代时，salutation被分配给切片中的下一个字符串值。
由于运行时调度器安排的goroutine可能会在将来的任何时间点运行，因此不确定在goroutine内将打印哪些值。 在我的机器上，在goroutines开始之前，循环很可能会退出。 
这意味着salutation变量超出了范围。 然后会发生什么？ goroutines仍然可以引用已经超出范围的东西吗？ 这个goroutine会访问可能已经被回收的内存吗？

这是关于Go如何管理内存的一个有趣的侧面说明。Go运行时足够敏锐地知道对salutation变量的引用仍然保留，因此会将内存传输到堆中，以便goroutine可以继续访问它。

在这个例子中，循环在任何goroutines开始运行之前退出，所以salutation转移到堆中，并保存对字符串切片“good day”中最后一个值的引用。所以会看到“good day”打印三次 。 编写该循环的正确方法是将salutation的副本传递给闭包，以便在运行goroutine时，它将对来自其循环迭代的数据进行操作：
```
var wg sync.WaitGroup
for _, salutation := range []string{"hello", "greetings", "good day"} {
    wg.Add(1)
    go func(salutation string) { // 1
        defer wg.Done()
        fmt.Println(salutation)
    }(salutation) // 2
}
wg.Wait()
```

在这里我们声明了一个参数，和其他的函数看起来差不多。我们将原始的salutation变量映射到更加明显的位置。
在这里，我们将当前迭代的变量传递给闭包。 一个字符串的副本被创建，从而确保当goroutine运行时，我们引用正确的字符串。
正如我们所看到的，我们得到的输出看起来没那么奇怪了：
```
good day 
hello 
greetings
```
这个例子的行为和我们预期的一样，只是稍微更冗长。多运行几次，输出顺序可能不同。
goroutine在相同的地址空间内运行，Go的编译器很好地处理了内存中的固定变量，因此goroutine不会意外地访问释放的内存，这允许开发人员专注于他们的问题而不是内存管理。
由于多个goroutine可以在相同的地址空间上运行，我们仍然需要担心同步问题。正如我们已经讨论过的，可以选择同步访问共享内存的例程访问，也可以使用CSP原语通过通信共享内存。
goroutines的另一个好处是它们非常轻巧。 这是官方FAQ的摘录：
新建立一个goroutine有几千字节，这样的大小几乎总是够用的。 如果出现不够用的情况，运行时会自动增加（并缩小）用于存储堆栈的内存，从而允许许多goroutine存在适量的内存中。
CPU开销平均每个函数调用大约三个廉价指令。 在相同的地址空间中创建数十万个goroutines是可以的。如果goroutines只是执行等同于线程的任务，那么系统资源的占用会更小。

每个goroutine几kb，那根本不是个事儿。让我们来亲手试着确认下。在此之前，我们必须了解一个关于goroutine的有趣的事：垃圾收集器不会收集以下形式的goroutines。如果我写出以下代码：
```go
go func() {
    // <操作会在这里永久阻塞>
}()
// Do work
```

这个goroutine将一直存在，直到整个程序退出。我们会在第四章的“防止Goroutine泄漏”中详细的聊一聊这个话题。
接下来，让我们回来看看该怎么写个例子来衡量一个goroutine的实际大小。
我们将goroutine不被垃圾收集的事实与运行时的自省能力结合起来，并测量在goroutine创建之前和之后分配的内存量：
```go
memConsumed := func() uint64 {
    runtime.GC()
    var s runtime.MemStats
    runtime.ReadMemStats(&s)
    return s.Sys
}

var c <-chan interface{}
var wg sync.WaitGroup
noop := func() { wg.Done(); <-c } // 1

const numGoroutines = 1e4 // 2
wg.Add(numGoroutines)
before := memConsumed() // 3
for i := numGoroutines; i > 0; i-- {
    go noop()
}
wg.Wait()
after := memConsumed() // 4
fmt.Printf("%.3fkb", float64(after-before)/numGoroutines/1000)
```

我们需要一个永不退出的goroutine，以便我们可以将它们中的一部分保存在内存中进行测量。 
不要担心我们目前如何实现这一目标。 只知道这个goroutine不会退出，直到这个过程结束。
这里我们定义要创建的goroutines的数量。 我们将使用大数定律渐近地逼近一个goroutine的大小。
这里测量创建分区之前所消耗的内存量。
这里测量创建goroutines后消耗的内存量。
在控制台会输出：
```
2.817kb
```

2017年7月这本书出版，go1.9发布于2017年8月24日，那么假设作者用的是当时最新的1.8版。译者用windows系统，go 1.10.1版，这个数字在8.908kb~9.186kb上下浮动。在centos6.4上测试，这个数字为2.748kb。

看起来文档是正确的。这个例子虽然有些理想化，但仍然让我们了解可能创建多少个goroutines有了大致的了解。

在我的笔记本上，我有8G内存，这意味着理论上我可以支持数百万的goroutines。当然，这忽略了在电脑上运行的其他东西。但这个快速估算的结果表明了goroutine是多么的轻量级。

存在一些可能会影响我们的goroutine规模的因素，例如上下文切换，即当某个并发进程承载的某些内容必须保存其状态以切换到其他进程时。如果我们有太多的并发进程，上下文切换可能花费所有的CPU时间，并且无法完成任何实际工作。在操作系统级别，使用线程，这样做代价可能会非常高昂。操作系统线程必须保存寄存器值，查找表和内存映射等内容，才能在操作成功后切换回当前线程。 然后它必须为传入线程加载相同的信息。

在软件中的上下文切换代价相对小得多。在软件定义的调度程序下，运行时可以更具选择性地进行持久检索，例如如何持久化以及何时发生持续化。我们来看看操作系统线程和goroutines之间上下文切换的相对性能。 首先，我们将利用Linux内置的基准测试套件来测量在同一内核的两个线程之间发送消息需要多长时间：
```
taskset -c 0 perf bench sched pipe -T
```
这会输出：
```shell
# Running 'sched/pipe' benchmark:
# Executed 1000000 pipe operations between two threads 

        Total time: 2.935 [sec]
            2.935784 usecs/op
               340624 ops/sec
```

这个基准测量实际上是衡量在一个线程上发送和接收消息所需的时间，所以我们将把结果分成两部分。 每个上下文切换1.467微秒。 这看起来不算太坏，但让我们先别急着下判断，再来比较下goroutine之间的上下文切换。

我们将使用Go构建一个类似的基准测试。下面的代码涉及到一些尚未讨论过的东西，所以如果有什么困惑的话，只需根据注释关注结果即可。 以下示例将创建两个goroutine并在它们之间发送消息：
```go
func BenchmarkContextSwitch(b *testing.B) {
    var wg sync.WaitGroup
    begin := make(chan struct{})
    c := make(chan struct{})

    var token struct{}
    sender := func() {
        defer wg.Done()
        <-begin //1
        for i := 0; i < b.N; i++ {
            c <- token //2
        }
    }
    receiver := func() {
        defer wg.Done()
        <-begin //1
        for i := 0; i < b.N; i++ {
            <-c //3
        }
    }

    wg.Add(2)
    go sender()
    go receiver()
    b.StartTimer() //4
    close(begin)   //5
    wg.Wait()
}
```

这里会被阻塞，直到接受到数据。我们不希望设置和启动goroutine影响上下文切换的度量。
在这里向接收者发送数据。struct{}{}是空结构体且不占用内存；这样我们就可以做到只测量发送信息所需要的时间。
在这里，我们接收传递过来的数据，但不做任何事。
开始启动计时器。
在这里我们通知发送和接收的goroutine启动。
我们运行该基准测试，指定只使用一个CPU，以便与之前的Linux基准测试想比较，我们来看看结果：
```go
go test -bench=. -cpu=1 /src/gos-concurrency-building-blocks/goroutines/fig-ctx-switch_test.go
BenchmarkContextSwitch  5000000 225ns/op
PASS        
ok  command-line-arguments  1.393s
```

每个上下文切换225 ns，哇！ 这是0.225μs，比我机器上的操作系统上下文切换快92％，如果你记得1.467μs的话。很难说有多少goroutines会导致过多的上下文切换，但我们可以很自然地说上限可能不会成为使用goroutines的障碍。

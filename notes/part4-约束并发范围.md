# 约束并发范围
    在使用并发代码时，安全操作有几种不同的选项。 我们已经使用并了解了其中两个：

    用于共享内存的同步原语(例如sync.Mutex)
    通过通信同步(例如channel)
    此外，有多个其他选项在多个并发进程中隐式安全：

        不可变数据
        受限制条件保护的数据
    从某种意义上讲，不可变数据是最理想的，因为它隐式地是并行安全的。每个并发进程可以在同一条的数据上运行，但不能修改它。如果要创建新数据，则必须创建所需修改数据的副本。 这不仅可以减轻开发人员认知负担，还可以让程序执行的更快(在某些情况下)。在Go中，可以通过使用值的副本而非该值的指针来实现此目的。 有些语言支持使用明确不变的值的指针； 然而，Go不在其中。

    不可变数据的使用依赖于约定——在我看来，坚持约定很难在任何规模的项目上进行协调，除非你有工具在每次有人提交代码时对代码进行静态分析。这里就有一个例子：

    data := make([]int, 4)

    loopData := func(handleData chan<- int) {
        defer close(handleData)
        for i := range data {
            handleData <- data[i]
        }
    }

    handleData := make(chan int)
    go loopData(handleData)

    for num := range handleData {
        fmt.Println(num)
    }
    我们可以看到，loopData函数和对handleData通道的循环都使用了整数切片data，但只有loopData对其进行了直接访问。
    但想想看，随着代码被其他的开发人员触及和修改，明显的，不明显的问题都有可能会被加入其中，并最终产生严重的错误（因为我们没有对data切片做显示的访问和操作约束）。正如我所提到的，一个静态分析工具可能会发现这类问题，但如此灵活的静态分析并不是很多团队能够实现的。 这就是为什么我更喜欢词汇约束，使用编译器来执行对变量的操作进行约束是非常好的。

    词法约束涉及使用词法作用域仅公开用于多个并发进程的正确数据和并发原语。 这使得做错事情变得不可能。 实际上，我们在第3章已经谈到了这个话题。回想一下通道部分，它讨论的只是将通道的读或写操作暴露给需要它们的并发进程。 我们再来看看这个例子：

    chanOwner := func() <-chan int {
        results := make(chan int, 5) //1
        go func() {
            defer close(results)
            for i := 0; i <= 5; i++ {
                results <- i
            }
        }()
        return results
    }

    consumer := func(results <-chan int) { //3
        for result := range results {
            fmt.Printf("Received: %d\n", result)
        }
        fmt.Println("Done receiving!")
    }

    results := chanOwner() //2
    consumer(results)
    这里我们在chanOwner函数的词法范围内实例化通道。这将导致通道的写入操作范围被限制在它下面定义的闭包中。 换句话说，它限制了这个通道的写入使用范围，以防止其他goroutine写入它。
    在这里，我们接受到一个只读通道，我们将它传递给消费者，消费者只能从中读取信息。
    这里我们收到一个int通道的只读副本。通过声明该函数的唯一用法是读取访问，我们将通道用法限制为只读。
    这样的设计方式就可以把通道的读取写入限制在一定的范围内。这个例子可能不是非常的有趣，因为通道是并发安全的。我们来看一个对非并发安全的数据结构约束的示例，它是一个bytes.Buffer实例：

    printData := func(wg *sync.WaitGroup, data []byte) {
        defer wg.Done()

        var buff bytes.Buffer
        for _, b := range data {
            fmt.Fprintf(&buff, "%c", b)
        }
        fmt.Println(buff.String())
    }

    var wg sync.WaitGroup
    wg.Add(2)
    data := []byte("golang")
    go printData(&wg, data[:3]) // 1
    go printData(&wg, data[3:]) // 2

    wg.Wait()
    这里我们传入包含前三个字节的data切片。
    这里我们传入包含剩余三个字节的data切片。
    在这个例子中，你可以看到，我们不需要通过通信同步内存访问或共享数据。

    那么这样做有什么意义呢？ 如果我们有同步功能，为什么要给予约束？ 答案是提高了性能并降低了开发人员的认知负担。同步带来了成本，如果你可以避免它，你就不必支付同步它们的成本。 你也可以通过同步回避所有可能的问题。利用词法约束的并发代码通常更易于理解。

    话虽如此，建立约束可能很困难，所以有时我们必须回到使用并发原语的开发思路上去。

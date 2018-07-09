# 关于chan
    Channel，即通道，衍生自Charles Antony Richard Hoare的CSP并发模型，是Go的并发原语，在Go语言中具有极其重要的地位。虽然它可用于同步内存的访问，但更适合用于goroutine之间传递信息。就像我们在之前“Go的并发哲学”章节所提到的那样，通道在任何规模的程序编码中都非常有用，因为它足够灵活，能够以各种方式组合在一起。我们在随后的“select语句”章节会进一步探讨其是如何构造的。

    可以把通道想象为一条河流，通道作为信息流的载体；数据可以沿着通道被传递，然后从下游被取出。基于这个原因，我通常用单词“Stream”为我的通道变量命名。在使用通道时，你把一个值传递给chan类型的变量，然后在程序的其他位置再把这个值从通道中读取出来。该值被传递的入口和出口不需要知道彼此的存在，你只需使用通道的引用进行操作就好。

    建立一个通道是非常简单的。下面这个例子展现了如何对通道声明和如何实例化。与Go中的其他类型一样，你可以使用 := 来简化通道的创建。不过在程序中通常会先声明通道，所以了解该如何将二者分成单独的步骤是很有用的：

    var dataStream chan interface{} // 1
    dataStream = make(chan interface{})  // 2
    1.这里声明了一个通道。我们说该通道的“类型”是interface{}的。
    2.这里我们使用内置函数make实例化通道。

    这个示例定义了一个名为dataStream的通道，在该通道上可以写入或读取任意类型的值（因为我们使用了空的接口）。通道也可以声明为仅支持单向数据流——即你可以定义仅支持发送或仅支持接收数据的通道。我将在本节的末尾解释单向数据流的重要性。

    要声明一个只能被读取的单向通道，你可以简单的使用 <- 符号，如下所示：

    var dataStream <-chan interface{}
    dataStream := make(<-chan interface{})
    与之相对应，要声明一个只能被发送的单向通道，把 <-放在 chan关键字的右侧即可：

    var dataStream chan<- interface{}
    dataStream := make(chan<- interface{})
    你不会经常看到实例化的单向通道，但你会经常看到它们被用作函数参数和返回类型，这是非常有用的，因为Go可以在需要时将双向通道隐式转换为单向通道，比如这样：

    var receiveChan <-chan interface{}
    var sendChan chan<- interface{}
    dataStream := make(chan interface{})

    // 这样做是有效的
    receiveChan = dataStream
    sendChan = dataStream
    要注意通道是有“类型”的。 在这个例子中，我们创建了一个接口“类型”的chan，这意味着我们可以在其上放置任何类型的数据，但是我们也可以给它一个更严格的类型来约束它可以传递的数据类型。 这是一个整数通道的例子：

    intStream := make(chan int)
    为了操作通道，我们再一次使用 <- 符号。我们看一个实际的例子：

    stringStream := make(chan string)
    go func() {
        stringStream <- "Hello channels!" //1
    }()
    fmt.Println(<-stringStream) //2
    这里我们将字符串放入通道stringStream。
    这里我们从通道中取出字符串并打印到标准输出流。
    这会输出：

    Hello channels!
    很简单，是吧。你所需要做的只是建立一个通道变量，然后将数据传递给它并从中读取数据。但是，尝试将值写入只读通道或从只写通道读取值都是错误的。如果我们尝试编译下面的例子，Go的编译器会报错：

    writeStream := make(chan<- interface{})
    readStream := make(<-chan interface{})

    <-writeStream
    readStream <- struct{}{}
    这会输出：

    invalid operation: <-writeStream (receive from send-only type chan<- interface {})
    invalid operation: readStream <- struct {} literal (send to receive-only type <-chan interface {})
    这是Go类型系统的一部分，即使在处理并发原语时也为我们保证类型安全。稍后我们会看到，这对构建易于推理的可组合逻辑程序提供了强大的保证。

    回想一下，在之前我们强调过，仅简单的定义一个goroutine并不能保证它在main goroutine退出之前运行。为此我们介绍了对sync包的各种使用案例。那么在使用通道的情况下该如何呢？看下面这个例子：

    该示例之所以产生这样的结果，是因为在Go中，通道是包含有阻塞机制的。这意味着试图写入已满的通道的任何goroutine都会等待直到通道被清空，并且任何尝试从空闲通道读取的goroutine都将等待，直到至少有一个元素被放置 。在这个例子中，我们的fmt.Println包含一个对通道stringStream的读取，并且将阻塞在那里，直到通道上被放置一个值。同样，匿名goroutine试图在stringStream上放置一个字符串，然后阻塞住等待被读取，所以goroutine在写入成功之前不会退出。因此，main goroutine和匿名的goroutine发生阻塞是毫无疑问的。

    stringStream := make(chan string)
    go func() {
        if 0 != 1 { //1
            return
        }
        stringStream <- "Hello channels!"
    }()

    fmt.Println(<-stringStream)
    在这里 我们确保通道stringStream永远不会获得值。
    这会产生错误：

    fatal error: all goroutines are asleep - deadlock!

    goroutine 1 [chan receive]: 
    main.main()

    /tmp/babel-23079IVB/go-src-230795Jc.go:15 +0x97 
    exit status 2
    main goroutine等待着stringSteam通道被放上一个值，而且由于我们的if条件，导致这不会发生。当匿名goroutine退出时，Go发现并报告死锁。在本节后面，我将解释如何构建我们的程序，以防止这种死锁。与此同时，让我们回到从通道读取数据的讨论。

    < - 运算符的接收形式也可以选择返回两个值，如下所示：

    stringStream := make(chan string)
    go func() {
        stringStream <- "Hello channels!"
    }()
    salutation, ok := <-stringStream //1
    fmt.Printf("(%v): %v", ok, salutation)
    我们在这里接收一个字符串salutation和一个布尔值ok。
    这会输出：

    (true): Hello channels!
    有意思吧。那么布尔值代表什么呢？这个值是读取操作的一个标识，用于指示读取的通道是由过程中其他位置的写入生成的值，还是由已关闭通道生成的默认值。 等一下; 一个已关闭的通道，那是什么？

    在程序中，能够指示还有没有更多值将通过通道发送是非常有用的。 这有助于下游流程知道何时移动，退出，或在新的通道上重新开启通信等。我们可以通过为每种类型提供特殊的标识符来完成此操作，但这会开发人员的工作产生巨大的重复性，如果能够内置将产生极大的便利，因此关闭通道就像是一个万能的哨兵，它说：“嘿，上游不会写更多的数据啦，做你想做的事吧。”要关闭频道，我们使用close关键字，就像这样：

    valueStream := make(chan interface{})
    close(valueStream)
    有趣的是，我们也可以从已关闭的通道读取。 看这个例子：

    intStream := make(chan int)
    close(intStream)
    integer, ok := <- intStream // 1
    fmt.Printf("(%v): %v", ok, integer)
    这里我们从已关闭的通道读取。
    这会输出：

    (false): 0
    注意我们在关闭通道前并没有把任何值放入通道。即便如此我们依然可以执行读取操作，而且尽管通道处在关闭状态，我们依然可以无限期地在此通道上执行读取操作。这是为了支持单个通道的上游写入器可以被多个下游读取器读取(在第四章我们会看到这是一种常见的情况)。第二个返回值——即布尔值ok——表明收到的值是int的零值，而非被放入流中传递过来。

    这为我们开辟了一些新的模式。首先是通道的range操作。与for语句一起使用的range关键字支持将通道作为参数，并且在通道关闭时自动结束循环。这允许对通道上的值进行简洁的迭代。 我们来看一个例子：

    intStream := make(chan int)
    go func() {
        defer close(intStream) // 1
        for i := 1; i <= 5; i++ {
            intStream <- i
        }
    }()

    for integer := range intStream { // 2
        fmt.Printf("%v ", integer)
    }
    在这里我们在通道退出之前保证正常关闭。这是一种很常见的Go惯用法。
    这里对intStream进行迭代。
    正如你所看到的，所有的值被打印后程序退出：

    1 2 3 4 5
    注意循环退出并没有设置条件，并且range也不返回第二个布尔值。对通道进行关闭的处理被隐藏了起来，以此保证循环的简洁。

    关闭某个通道同样可以被作为向多个goroutine同时发生消息的方式之一。如果你有多个goroutine在单个通道上等待，你可以简单的关闭通道，而不是循环解除每一个goroutine的阻塞。由于一个已关闭的通道可以被无限次的读取，因此其中有多少goroutine在阻塞状态并不重要，关闭通道(以解除所有阻塞)消耗的资源又少执行的速度又快。以下是一次解除多个goroutine的示例：

    begin := make(chan interface{})
    var wg sync.WaitGroup
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            <-begin //1
            fmt.Printf("%v has begun\n", i)
        }(i)
    }

    fmt.Println("Unblocking goroutines...")
    close(begin) //2
    wg.Wait()
    这里对begin通道进行读取，由于通道中没有任何值，会产生阻塞。
    这里我们关闭通道，这样所有goroutine的阻塞会被解除。
    你可以看到，在我们关闭begin 通道之前，没有任何一个goroutine开始运行：

    Unblocking goroutines...
    4 has begun
    2 has begun
    3 has begun
    0 has begun
    1 has begun
    回想一下在“sync包”中我们讨论过sync.Cond实现类似功能的例子，你当然可以使用Single或者Brocast来做，不过通道是可组合的，所以这也是我最喜欢的同时解除多个goroutine阻塞的方法。

    接下来，我们来讨论“缓冲通道”，这种通道实在实例化时候提供可携带元素的容量。这意味着，即使没有对通道进行读取操作，goroutine仍然可以执行n次写入，这里的n即缓冲通道的容量。下面是一个实例化的例子：

    var dataStream chan interface{} 
    dataStream = make(chan interface{}, 4) 
    这里我们创建一个容量为4的缓冲通道。 这意味着我们可以将4个元素放在通道上，而不管它是否被读取（在数量达到上限之前，写入行为都不会发生阻塞）。
    我们再一次把初始化分为了两行，这样你可以清楚的发现，一个缓冲通道和一个非缓冲通道
    在声明上上没有区别的(区别只在实例化部分)。有趣的地方在于，我们可以在实例化的位置对通道是否是缓冲的进行控制。这表明，通道的建立应该与goroutine紧密结合，这样我们可以极大的提高代码的可读性。

    无缓冲的通道也可以按缓冲通道定义：无缓冲的通道可以视作一个容量为0的缓冲通道。就像下面这样：

    a := make(chan int)
    b := make(chan int, 0)
    这两个通道都是int“类型”的。请记住我们在讨论“阻塞”时所代表的含义，我们说向一个已满的通道写入，会出现阻塞，从一个已空的通道读取，也会出现阻塞。这里的“满”和“空”是针对容量或缓冲区大小而言的。无缓冲的通道所拥有的容量为0，所以任何写入行为之后它都会是满的。一个容量为4的缓冲通道在4次写入后会是满的，并且会在第5次写入时出现阻塞，因为它已经没有其他位置可以放置第5个元素，这时它表现出的行为与无缓冲通道一样：由此可见，缓冲通道和无缓冲通道的区别在于，通道为空和满的前提条件是不同的。通过这种方式，缓冲通道可以在内存中构建用于并发进程通信的FIFO队列。

    为了帮助理解这一点，我们来举例说明缓冲通道容量为4的示例中发生了什么。 首先，让我们初始化它：

    c := make(chan rune, 4)
    从逻辑上讲，这会创建一个带有四个空位的缓冲区的通道：



    现在，让我们向通道写入：

    c <- 'A'
    当这个通道没有被读取时，A字符将被放置在通道缓冲区的第一个空位中，像这样：



    随后每次写入缓冲通道(同样假设没有被读取)，将填充缓冲通道中的剩余空位，像这样：

    c <- 'B'


    c <- 'C'


    c <- 'D'


    经过四次写入，我们的缓冲通道已经装满了4个元素。如果我们再向通道中进行写入的话：

    c <- 'E'


    当前的goroutine会表现为阻塞！并且goroutine将一直保持阻塞状态，直到由其他的goroutine执行读取操作在缓冲区中空出了位置。 让我们看看是什么样子的：

    <- c


    正如你所看到的那样，读取时会接收到放在通道上的第一个字符A，被阻塞的写入阻塞解除，E被放置在缓冲区的末尾。

    如果，如果缓冲通道为空且有接收器读取，则缓冲器将被绕过，并且该值将直接从发送器传递到接收器（存疑）。实际上，这是透明地发生的，但值得了解。

    缓冲通道在某些情况下很有用，但你应该小心使用。正如我们将在下一章中看到的那样，缓冲通道很容易成为不成熟的优化，并且通过使用它们死锁会变得更为隐蔽。我猜你宁愿在第一次编写代码时找到一个死锁，而不是在半夜系统停机的时候。

    让我们来看看另一个更完整的代码示例，以便更好地了解缓冲通道的工作方式：

    var stdoutBuff bytes.Buffer         //1
    defer stdoutBuff.WriteTo(os.Stdout) //2

    intStream := make(chan int, 4) //3
    go func() {
        defer close(intStream)
        defer fmt.Fprintln(&stdoutBuff, "Producer Done.")
        for i := 0; i < 5; i++ {
            fmt.Fprintf(&stdoutBuff, "Sending: %d\n", i)
            intStream <- i
        }
    }()

    for integer := range intStream {
        fmt.Fprintf(&stdoutBuff, "Received %v.\n", integer)
    }
    这里我们创建一个内存缓冲区来帮助缓解输出的不确定性。 它不会给带来我们任何保证，但比直接写stdout要快一些。
    在这里，我们确保在进程退出之前将缓冲区内容写入标准输出。
    这里我们创建一个容量为4的缓冲通道。
    在这个例子中，写入stdout的顺序是不确定的，但你仍然可以大致了解匿名goroutine是如何工作的。 如果你检查输出结果，可以看到我们的匿名goroutine能够将所有五个结果放在intStream中，并在主要goroutine将一个结果关闭之前退出。

    Sending: 0
    Sending: 1
    Sending: 2
    Sending: 3
    Sending: 4
    Producer Done.
    Received 0 
    Received 1 
    Received 2 
    Received 3 
    Received 4 
    这是一个在正确条件下可以使用的优化示例：如果写入通道的goroutine明确知道将会写入多少信息，则创建相对应的缓冲通道容量会很有用，就可以尽可能快地进行读取。 当然，这样做是有限制的，我们将在下一章中介绍。

    我们已经讨论了无缓冲的频道，缓冲频道，双向频道和单向频道。目前还没有讨论到的还有通道的默认值：nil。程序是如何处理处理nil通道的呢？首先，让我们试着从一个nil通道中读取：

    var dataStream chan interface{}
    <-dataStream
    这会输出：

    fatal error: all goroutines are asleep - deadlock!

    goroutine 1 [chan receive (nil chan)]:
    main.main()
        F:/code/gospcace/src/myConcurrency/l1introduction/l01/main.go:6 +0x30
    死锁出现了。这说明从一个nil通道进行读取会阻塞程序(注意，这段代码的前提是在main函数中执行，所以会导致死锁。如果是运行在单个gouroutine中，那么就不会是死锁而是阻塞)。让我们再试试写入：

    var dataStream chan interface{}
    dataStream <- struct{}{}
    这会输出：

    fatal error: all goroutines are asleep - deadlock!

    goroutine 1 [chan send (nil chan)]:
    main.main()
        F:/code/gospcace/src/myConcurrency/l1introduction/l01/main.go:6 +0x53
    看来对一个nil通道进行写入操作同样会阻塞。
    我们再试试关闭操作：

    var dataStream chan interface{}
    close(dataStream)
    这会输出：

    panic: close of nil channel

    goroutine 1 [running]:
    main.main()
        F:/code/gospcace/src/myConcurrency/l1introduction/l01/main.go:6 +0x31
    程序挂掉了，这也许是最符合你预期的结果。无论如何，请务必确保你的通道在工作前已经完成了初始化。

    我们已经了解了很多与通道互动的规则。现在你已经了解了在通道上执行操作的方式和为什么这样做的原因。 下表列举了通道上的操作对应状态的通道会发生什么。

    注意：表中的用词都很简短，为了减少不必要的歧义或混乱，并未对该表进行不必要的翻译，此外，正如上面例子所展现的，该表的操作结果默认都是在main函数下操作。请以批判的眼光审视下表。



    如果我们查看该表，可以察觉到在操作中可能产生问题的地方。这里有三个可能导致阻塞的操作，以及三个可能导致程序恐慌的操作。乍看之下，通道的使用上限制很多，但在检查了这个限制产生的动机并熟悉了通道的使用后，它变得不那么可怕并开始具有很大意义。让我们讨论如何组织不同类型的通道来构筑稳健的程序。

    我们应该做的第一件事是将通道置于正确的环境中，即分配通道所有权。我将所有权定义为goroutine的实例化，写入和关闭。就像在那些没有垃圾回收的语言中使用内存一样，重要的是要明确哪个goroutine拥有该通道，以便从逻辑上推理我们的程序。单向通道声明是一种工具，它可以让我们区分哪些gouroutine拥有通道，哪些goroutine仅使用通道：通道所有者对通道具有写入访问权（chan或chan<- ），而通道使用者仅具有读取权（<-chan）。一旦我们对通道权责区分，上表的结果自然就会出现。我们可以开始对拥有通道和不拥有通道的goroutine赋予不同的责任并给予对应的检查(以增强程序和逻辑的健壮性)。

    让我们从通道的所有者说起。当一个goroutine拥有一个通道时应该：

    初始化该通道。
    执行写入操作，或将所有权交给另一个goroutine。
    关闭该通道。
    将此前列出的三件事情封装在一个列表中，并通过订阅通道将其公开。
    通过将这些责任分配给通道所有者，会发生一些事情：

    因为我们是初始化频道的人，所以我们要了解写入nil通道会带来死锁的风险。
    因为我们是初始化频道的人，所以我们要了解关闭ni通道会带来恐慌的风险。
    因为我们是决定频道何时关闭的人，所以我们要了解写入已关闭的通道会带来恐慌的风险。
    因为我们是决定何时关闭频道的人，所以我们要了解多次关闭通道会带来恐慌的风险。
    我们在编译时使用类型检查器来防止对通道进行不正确的写入。
    现在我们来看看读取时可能发生的那些阻塞操作。 作为一个通道的消费者，我只需要担心两件事情：

    通道什么时候会被关闭。
    处理基于任何原因出现的阻塞。
    解决第一个问题，通过检查读取操作的第二个返回值就可以。第二点很难，因为它取决于你的算法(和业务逻辑)：你可能想要超时，当获得通知时你可能想停止读取操作，或者你可能只是满足于在整个生命周期中产生阻塞。 重要的是，作为一个消费者，你应该明确这样一个事实，即读取操作可以并且必将产生阻塞。我们将在下一章中探讨如何实现select语句解决这个棘手的问题。

    现在，让我们用一个例子来总结以上的思考结果。我们建立一个goroutine，它拥有一个通道，一个消费者，它会处理阻塞问题：

    chanOwner := func() <-chan int {

        resultStream := make(chan int, 5)//1
        go func() {//2
            defer close(resultStream)//3
            for i := 0; i <= 5; i++ {
                resultStream <- i
            }
        }()
        return resultStream//4

    }

    resultStream := chanOwner()
    for result := range resultStream {//5
        fmt.Printf("Received: %d\n", result)
    }
    fmt.Println("Done receiving!")
    这里我们实例化一个缓冲通道。 由于我们知道我们将产生六个结果，因此我们创建了五个缓冲通道，以便该goroutine可以尽快完成操作。
    在这里，我们启动一个匿名的goroutine，它在resultStream上执行写操作。 请注意，我们是如果创建goroutines的， 它现在被封装在函数中。
    这里我们确保resultStream在操作完成后关闭。作为通道所有者，这是我们的责任。
    我们在这里返回通道。由于返回值被声明为只读通道，resultStream将隐式转换为只读的。
    这里我们消费了resultStream。 作为消费者，我们只关心阻塞和通道的关闭。
    这会输出：

    Received: 0
    Received: 1
    Received: 2
    Received: 3
    Received: 4
    Received: 5
    Done receiving!
    注意resultStream通道的生命周期如何封装在chanOwner函数中。很明显，写入不会发生在nil或已关闭的频道上，并且关闭总是会发生一次。这消除了我们之前提到的部分风险。我强烈建议你在自己的程序中尽可能做到保持通道覆盖范围最小，以便这些事情保持明显。如果你将一个通道作为一个结构体的成员变量，并且有很多方法，它很快就会把你自己给绕进去（虽然很多库和书中都这么干，但只有这本书的作者将这一点给明确提出来了）。

    消费者功能只能读取通道，因此只需知道应该如何处理阻塞读取和通道关闭。 在这个小例子中，我们采取了这样的方式：在通道关闭之前阻塞程序是完全没问题的。

    如果你设计自己的代码时来遵循这个原则，那么对你的系统进行推理就会容易得多，而且它很可能会像你期望的那样执行。我不能保证你永远不会引入阻塞或恐慌，但是当你这样遇到这样的情况时，我认为你会发现你的通道所有权范围要么太大，要么所有权不清晰。

    通道是首先吸引我使用Go的原因之一。 结合goroutines和闭包的简约性，编写干净、正确的并发代码是比较容易的。在很多方面，通道是将goroutine绑在一起的胶水。 本节为你概述了什么是通道以及如何使用它们。当我们开始编写通道以形成更高阶的并发设计模式时，真正的乐趣就开始了。我们将在下一章中体会到这一点。

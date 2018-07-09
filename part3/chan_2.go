package main

import "fmt"

func main() {
	//生产者返回只读通道
	chanOwner := func() <-chan int {
		resultStream := make(chan int, 5) //1
		go func() {                       //2
			defer close(resultStream) //3
			for i := 0; i <= 5; i++ {
				resultStream <- i
			}
		}()
		return resultStream //4

	}

	resultStream := chanOwner()
	//消费者读取通道中的内容
	for result := range resultStream { //5
		fmt.Printf("Received: %d\n", result)
	}
	fmt.Println("Done receiving!")
}

/**
1 这里我们实例化一个缓冲通道。 由于我们知道我们将产生六个结果，因此我们创建了五个缓冲通道，以便该goroutine可以尽快完成操作。
2 在这里，我们启动一个匿名的goroutine，它在resultStream上执行写操作。 请注意，我们是如果创建goroutines的， 它现在被封装在函数中。
3 这里我们确保resultStream在操作完成后关闭。作为通道所有者，这是我们的责任。
4 我们在这里返回通道。由于返回值被声明为只读通道，resultStream将隐式转换为只读的。
5 这里我们消费了resultStream。 作为消费者，我们只关心阻塞和通道的关闭。
*/

package main

import (
	"fmt"
	"log"
	"runtime"
)

func produce(p chan<- int) {
	for i := 0; i < 10; i++ {
		runtime.Gosched()
		p <- i
		fmt.Println("send:", i)
	}
}
func consumer(c <-chan int, done chan struct{}) {
	for i := 0; i < 10; i++ {
		v := <-c
		fmt.Println("receive data:", v)
	}

	close(done) //关闭通道后,<-done就会立即释放
}

func main() {
	ch := make(chan int)
	done := make(chan struct{}) //通过完成信号量实现通道之间通信
	go produce(ch)
	go consumer(ch, done)
	// time.Sleep(1 * time.Second) 不建议这么处理
	<-done

	ch2 := make(chan int)
	go produce(ch2) //生产者
	//对数据进行消费,从通道中取出数据
	//当生产者没有准备好,就会阻塞,直到生产者发送了数据到通道中,消费者才可以读取到数据
	//双方都需要准备好
	for i := 0; i < 10; i++ {
		log.Println("获取到的chan: ", <-ch2)
	}

	//缓冲通道
	ch3 := make(chan int, 10)

	//下面的生产者,可以在独立的携程中处理,也可以不在携程中处理
	go func() {
		for i := 0; i < 10; i++ {
			runtime.Gosched()
			ch3 <- i
			fmt.Println("has send:", i)
		}

		close(ch3) //这里必须关闭,否则for...range会一直阻塞
	}()

	for value := range ch3 {
		log.Println("value: ", value)
	}
}

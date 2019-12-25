package main

import "log"

func main() {

	//通过chan实现同步
	//ch相当于同步信号量
	ch := make(chan struct{})
	go func() {
		defer log.Println("子携程将要退出")

		for i := 0; i < 10; i++ {
			log.Println("current index: ", i)
		}

		ch <- struct{}{}
	}()

	<-ch

	// hello()
	// hello2()

	// hello3() // 产生死锁

	log.Println("main will exit")
}

func hello() {
	ch2 := make(chan int, 3)
	// 缓冲区大小3，当缓冲区满了的时候，生产者就无法放入
	// 这个时候需要消费者把数据从缓冲区通道中，把数据拿走，当消费者拿走了一部分后
	// 生产者那一方，又可以放入数据进来

	//异步生产者
	go func() {
		for i := 0; i < 10; i++ {
			ch2 <- i
		}
	}()

	// 消费者
	for i := 0; i < 10; i++ {
		log.Println("ch2 index: ", <-ch2)
	}
}

func hello2() {
	ch2 := make(chan int, 3)
	// 异步消费者
	go func() {
		for i := 0; i < 10; i++ {
			log.Println("ch2 index: ", <-ch2)
		}
	}()

	for i := 0; i < 10; i++ {
		ch2 <- i
	}

}

// fatal error: all goroutines are asleep - deadlock，产生死锁
func hello3() {
	ch2 := make(chan int, 3)
	for i := 0; i < 10; i++ {
		ch2 <- i
	}

	for i := 0; i < 10; i++ {
		log.Println("ch2 index: ", <-ch2)
	}
}

/**
$ go run chan_demo.go
2019/12/25 23:01:51 current index:  0
2019/12/25 23:01:51 current index:  1
2019/12/25 23:01:51 current index:  2
2019/12/25 23:01:51 current index:  3
2019/12/25 23:01:51 current index:  4
2019/12/25 23:01:51 current index:  5
2019/12/25 23:01:51 current index:  6
2019/12/25 23:01:51 current index:  7
2019/12/25 23:01:51 current index:  8
2019/12/25 23:01:51 current index:  9
2019/12/25 23:01:51 子携程将要退出
2019/12/25 23:01:51 main will exit
*/


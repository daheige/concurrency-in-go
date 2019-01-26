/**
====关于chan在通信过程中，生产者和消费者运行模式研究===
采用无缓冲通道：需要一方准备好姿势
	让生产者和消费者其中的一个在独立协程中运行
	1. go product....consumer
	2. go Consumer...product 流程上异步消费者要先加入主携程的一条线上,
							 当取不到数据，就会一直等待生产者发送数据过来
	让生产者和消费者都在独立携程中处理
	3. go product...go consumer 借助空结构体chan信号量告诉消费者可以开始消费任务了
	4. product...consumer 这种模式，携程阻塞，陷入死锁，无法运行,这是一种错误的运行模式
采用有缓冲通道: 无需准备好姿势
	1. go product...consumer
	2. product...go consumer
	3. product...consumer
	4. go product...go consumer 需要借助信号量阻塞主携程退出或for{ runtime.Gosched() }模式不让主携程退出
								这种模式，不推荐使用，cpu会频繁切换上下文，资源消耗比较大
*/
package main

import (
	"log"
)

func main() {
	//===========无缓冲通道的生产者和消费者============
	ch := make(chan int)
	//模式1： 生产者异步--消费者在主协程中
	//go product(ch)
	//consumer(ch) //消费者取出任务，当通道中没有数据，会一直阻塞，直到有数据发送，消费者就会接收

	//模式2： 生产者同步发送---消费者异步消费
	// 流程上异步消费者要先加入主携程的一条线上
	//go consumer(ch)
	//product(ch)

	//done := make(chan struct{})

	//下面生产者同步执行，消费者异步处理，会抛出异常
	//fatal error: all goroutines are asleep - deadlock!
	/* product(ch)
	go consumer(ch)*/

	//模式3： 生产者和消费者都独立携程处理,借助空结构体chan信号量告诉消费者可以开始消费任务了
	done := make(chan struct{})
	go productDone(ch, done)

	go consumer(ch)
	<-done //采用通道信号量保证独立携程都可以运行完毕

	log.Println("====end====")

	//===========有缓冲通道的生产者和消费者============
	log.Println("====ch2===")
	ch2 := make(chan int, 10)
	go product(ch2)
	consumer(ch2)

	log.Println("====ch3===")
	ch3 := make(chan int, 10)
	product(ch3)
	go consumer(ch3)

	log.Println("====ch4===")
	ch4 := make(chan int, 10)
	product(ch4)
	consumer(ch4)

	log.Println("exit")

	//在有缓冲通道模式下，如果消费者和生产者都在独立携程中跑，程序在主携程退出后，就会退出
	//除非采用信号量阻塞主携程不让它退出，或for{ runtime.Gosched() }模式不让主携程退出
	//这种模式，不推荐使用
	/*log.Println("====ch5===")
	ch5 := make(chan int, 10)
	go product(ch5)
	go consumer(ch5)

	for {
		runtime.Gosched()
	}*/
}

func product(ch chan<- int) {
	defer close(ch)
	for i := 0; i < 10; i++ {
		log.Println("has send: ", i)
		ch <- i
	}

	log.Println("has send finish")
}

func consumer(ch <-chan int) {
	for i := 0; i < 10; i++ {
		log.Println("has recive: ", <-ch)
	}

	log.Println("has cust success")
}

func productDone(ch chan<- int, done chan<- struct{}) {
	defer close(ch)
	defer close(done) //关闭通道
	for i := 0; i < 10; i++ {
		log.Println("has send: ", i)
		ch <- i
	}

	log.Println("has send finish")
}

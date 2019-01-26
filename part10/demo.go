/**
采用无缓冲通道实现生产者异步携程处理，消费者在主携程中消费
 */
package main

import "log"

func main() {
	ch := make(chan int)
	//生产者是异步执行,如果是非go func模式的话，就会陷入死锁
	// 没有go就会报错,fatal error: all goroutines are asleep - deadlock!
	go product1(ch)

	consumer1(ch) //消费者取出任务，当通道中没有数据，会一直阻塞，直到有数据发送，消费者就会接收
	log.Println("exit")
}

func product1(ch chan<- int){
	defer close(ch) //生产者关闭通道，停止发送数据
	for i := 0;i<10;i++{
		log.Println("has send: ",i)
		ch <- i
	}

	log.Println("has send finish")
}

func consumer1(ch <-chan int){
	for i := 0;i<10;i++{
		log.Println("has recive: ",<-ch)
	}

	log.Println("has cust success")
}

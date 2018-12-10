package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan struct{}) //用来控制流程,如果是用作通知一般建议用空struct不占用空间
	c := make(chan int)
	go func() {
		defer close(done) //关闭通道
		for x := range c {
			fmt.Println("接收到c", x)
		}
	}()

	//发送方
	for i := 0; i < 4; i++ {
		c <- i
	}

	close(c)
	<-done //当done没有值,会一直阻塞,当done有值或done通道被关闭,才会释放

	//采用空结构体实现阻塞,当关闭通道后,立即释放
	exit := make(chan struct{})
	go func() {
		// defer close(exit)
		time.Sleep(1 * time.Second)
		fmt.Println("goroutine has done")
		close(exit)
	}()

	<-exit //通道关闭后,立即释放
	fmt.Println("task has done")
}

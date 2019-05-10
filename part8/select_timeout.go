package main

import (
	"fmt"
	"time"
)

func main() {
	var ch = make(chan int)
	quit := make(chan bool) //退出独立goroutine信号
	go func() {
		for {
			select {
			case num := <-ch:
				fmt.Println("num is ", num)
			case <-time.After(3 * time.Second):
				fmt.Println("timeout")
				quit <- true
			}
		}
	}()

	//模拟超时
	for i := 0; i < 5; i++ {
		fmt.Println("i = ", i)
		ch <- i
		time.Sleep(600 * time.Millisecond)
	}

	<-quit
	fmt.Println("end")

}

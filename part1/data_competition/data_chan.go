package main

import (
	"fmt"
)

//利用chan保证goroutine执行完毕
func main() {
	var data int

	ch := make(chan bool) //通过chan通道来保证goroutine执行
	go func() {
		data++
		ch <- true //这里发送value到ch中
	}()

	//取出通道中的值
	v := <-ch //这里会阻塞,一直等待ch的发送者把数据放入ch中,从而保证goroutine执行完毕
	fmt.Println(v)
	if data == 1 {
		fmt.Printf("the value is %v.\n", data)
	}

	fmt.Println(111)
}

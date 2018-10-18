package main

import "fmt"

func main() {
	//通过chan实现goroutine同步执行
	var ch = make(chan struct{})
	go func() {
		defer fmt.Println("goroutine执行完毕")
		for i := 0; i < 10; i++ {
			fmt.Println("i=", i)
		}
		close(ch)
	}()

	<-ch //当通道中没有值的话，一直阻塞，直到接受到值(或者通道关闭)为止
}

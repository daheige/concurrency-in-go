package main

import "fmt"

func main() {
	var ch = make(chan int) //双向通道
	go production(ch)
	cust(ch)
}

//单方向的通道,写入
func production(ch chan<- int) {
	for i := 0; i < 3; i++ {
		ch <- i
	}

	//写入后，关闭通道
	close(ch) //这里必须关闭，range读取的时候，才会释放
}

//只能读取
func cust(ch <-chan int) {
	for num := range ch {
		fmt.Println("read num: ", num)
	}
}

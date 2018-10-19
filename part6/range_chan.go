package main

import "fmt"

func main() {
	var ch = make(chan int, 10)
	for i := 0; i < 10; i++ {
		fmt.Println("开始放入", i)
		ch <- i
	}
	close(ch) //缓冲满了后，就关闭通道

	for num := range ch {
		fmt.Println("读取num: ", num)
	}
}

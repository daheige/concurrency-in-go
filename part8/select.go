package main

import "fmt"

func main() {
	x := makeRndNum()
	fmt.Println(x)
}

func makeRndNum() int {
	x := make(chan int)
	//独立goroutine生成0,1
	go func() {
		select { //随机选择一个chan写入
		case x <- 1:
			fmt.Println("x = 1")
		case x <- 0:
			fmt.Println("x = 0")
		}
	}()

	return <-x
}
